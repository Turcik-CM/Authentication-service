package postgres

import (
	"auth-service/pkg/hashing"
	"auth-service/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"

	pb "auth-service/genproto/user"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) storage.UserStorage {
	return &UserRepo{db: db}
}

func (p *UserRepo) Create(req *pb.CreateRequest) (*pb.UserResponse, error) {
	userID := uuid.New().String()

	query := `INSERT INTO users (id, phone, email, password) VALUES ($1, $2, $3, $4) RETURNING id`
	_, err := p.db.Exec(query, userID, req.Phone, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	profileQuery := `INSERT INTO user_profile (user_id, first_name, last_name, username, nationality, bio) 
	                 VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = p.db.Exec(profileQuery, userID, req.FirstName, req.LastName, req.Username, req.Nationality, req.Bio)
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{Id: userID, Email: req.Email, FirstName: req.FirstName, LastName: req.LastName}, nil
}

func (p *UserRepo) GetProfile(req *pb.Id) (*pb.GetProfileResponse, error) {
	query := `SELECT u.id, u.email, u.phone, p.first_name, p.last_name, p.username, p.nationality, p.bio, 
	                   p.followers_count, p.following_count, p.posts_count
	          FROM users u
	          JOIN user_profile p ON u.id = p.user_id
	          WHERE u.id = $1 and p.role != 'admin' and u.deleted_at = 0`

	row := p.db.QueryRow(query, req.UserId)
	var res pb.GetProfileResponse
	err := row.Scan(&res.UserId, &res.Email, &res.PhoneNumber, &res.FirstName, &res.LastName, &res.Username, &res.Nationality,
		&res.Bio, &res.FollowersCount, &res.FollowingCount, &res.PostsCount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &res, nil
}

func (p *UserRepo) UpdateProfile(req *pb.UpdateProfileRequest) (*pb.UserResponse, error) {
	query := `UPDATE user_profile 
	          SET first_name = $1, last_name = $2, username = $3, nationality = $4, bio = $5, profile_image = $6, updated_at = now()
	          WHERE user_id = $7 RETURNING user_id`

	_, err := p.db.Exec(query, req.FirstName, req.LastName, req.Username, req.Nationality, req.Bio, req.ProfileImage, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{Id: req.UserId, FirstName: req.FirstName, LastName: req.LastName, Email: ""}, nil
}

func (p *UserRepo) ChangePassword(req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	query := `SELECT password FROM users WHERE id = $1 and deleted_at = 0`

	row := p.db.QueryRow(query, req.UserId)
	var password string
	err := row.Scan(&password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	ok := hashing.CheckPasswordHash(password, req.CurrentPassword)
	if !ok {
		return nil, errors.New("password is incorrect")
	}
	query = `UPDATE users SET password = $1, updated_at = now() WHERE id = $2`
	_, err = p.db.Exec(query, req.NewPassword, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.ChangePasswordResponse{Message: "Password updated successfully"}, nil
}

func (p *UserRepo) ChangeProfileImage(req *pb.URL) (*pb.Void, error) {
	query := `UPDATE user_profile SET profile_image = $1, updated_at = now() WHERE user_id = $2`
	_, err := p.db.Exec(query, req.Url, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.Void{}, nil
}

func (p *UserRepo) FetchUsers(req *pb.Filter) (*pb.UserResponses, error) {
	query := `SELECT u.id, u.email, p.first_name, p.last_name, p.username, u.created_at
	          FROM users u
	          JOIN user_profile p ON u.id = p.user_id
	          WHERE p.username ILIKE $1 and p.role = 'user'
	          LIMIT $2 OFFSET $3`

	rows, err := p.db.Query(query, req.FirstName, req.Limit, (req.Page-1)*req.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*pb.UserResponse
	for rows.Next() {
		var user pb.UserResponse
		if err := rows.Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.Username, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return &pb.UserResponses{Users: users}, nil
}

func (p *UserRepo) ListOfFollowing(req *pb.Id) (*pb.Followings, error) {
	followings := &pb.Followings{}

	query := `
		SELECT p.username
		FROM follows f
		JOIN user_profile p ON f.following_id = p.user_id
		WHERE f.follower_id = $1;
    `

	rows, err := p.db.Query(query, req.UserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var followingID string
		if err := rows.Scan(&followingID); err != nil {
			return nil, err
		}
		followings.Following = append(followings.Following, followingID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return followings, nil
}

func (p *UserRepo) ListOfFollowingByUsername(req *pb.Id) (*pb.Followings, error) {
	followings := &pb.Followings{}

	query := `
		SELECT p2.username AS following_username
		FROM follows f
		JOIN user_profile p1 ON f.follower_id = p1.user_id
		JOIN user_profile p2 ON f.following_id = p2.user_id
		WHERE p1.username = $1;
    `

	rows, err := p.db.Query(query, req.UserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var followingUsername string
		if err := rows.Scan(&followingUsername); err != nil {
			return nil, err
		}
		followings.Following = append(followings.Following, followingUsername)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return followings, nil
}

func (p *UserRepo) ListOfFollowersByUsername(req *pb.Id) (*pb.Followers, error) {
	followers := &pb.Followers{}
	query := `
		SELECT p2.username AS follower_username
		FROM follows f
		JOIN user_profile p1 ON f.following_id = p1.user_id
		JOIN user_profile p2 ON f.follower_id = p2.user_id
		WHERE p1.username = $1;
    `

	rows, err := p.db.Query(query, req.UserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var followerUsername string
		if err := rows.Scan(&followerUsername); err != nil {
			return nil, err
		}
		followers.Followers = append(followers.Followers, followerUsername)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return followers, nil
}

func (p *UserRepo) ListOfFollowers(req *pb.Id) (*pb.Followers, error) {
	followers := &pb.Followers{}
	query := `
		SELECT p.username
		FROM follows f
		JOIN user_profile p ON f.follower_id = p.user_id
		WHERE f.following_id = $1;
    `

	rows, err := p.db.Query(query, req.UserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var followerID string
		if err := rows.Scan(&followerID); err != nil {
			return nil, err
		}
		followers.Followers = append(followers.Followers, followerID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return followers, nil
}

func (p *UserRepo) DeleteUser(req *pb.Id) (*pb.Void, error) {
	query := `UPDATE users SET deleted_at = EXTRACT(EPOCH FROM NOW()) WHERE id = $1 AND deleted_at = 0`

	_, err := p.db.Exec(query, req.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to mark user as deleted: %w", err)
	}

	return &pb.Void{}, nil
}

// ----------------------------------------------------------------------------------------

func (p *UserRepo) Follow(in *pb.FollowReq) (*pb.FollowRes, error) {
	query := `INSERT INTO follows (follower_id, following_id, created_at)
	          VALUES ($1, $2, NOW())
	          RETURNING follower_id, following_id, created_at`

	var res pb.FollowRes
	err := p.db.QueryRowContext(context.Background(), query, in.FollowerId, in.FollowingId).Scan(
		&res.FollowerId, &res.FollowingId, &res.FollowedAt)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (p *UserRepo) Unfollow(in *pb.FollowReq) (*pb.DFollowRes, error) {
	query := `DELETE FROM follows WHERE follower_id = $1 AND following_id = $2
	          RETURNING follower_id, following_id`

	var res pb.DFollowRes
	err := p.db.QueryRowContext(context.Background(), query, in.FollowerId, in.FollowingId).Scan(
		&res.FollowerId, &res.FollowingId)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no such follow relation exists")
		}
		return nil, err
	}

	return &res, nil
}

func (p *UserRepo) GetUserFollowers(in *pb.Id) (*pb.Count, error) {
	query := `SELECT COUNT(*) FROM follows WHERE following_id = $1`

	var count pb.Count
	err := p.db.QueryRowContext(context.Background(), query, in.UserId).Scan(&count.Count)
	if err != nil {
		return nil, err
	}

	return &count, nil
}

func (p *UserRepo) GetUserFollows(in *pb.Id) (*pb.Count, error) {
	query := `SELECT COUNT(*) FROM follows WHERE follower_id = $1`

	var count pb.Count
	err := p.db.QueryRowContext(context.Background(), query, in.UserId).Scan(&count.Count)
	if err != nil {
		return nil, err
	}

	return &count, nil
}

func (p *UserRepo) MostPopularUser(in *pb.Void) (*pb.UserResponse, error) {
	var (
		userID string
	)

	query := `SELECT follower_id FROM follows ORDER BY follower_id DESC LIMIT 1`

	var user pb.UserResponse
	err := p.db.QueryRow(query).Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get most popular user: %w", err)
	}

	query1 := `SELECT u.id, u.email, u.phone, p.first_name, p.last_name, p.username, p.nationality, p.bio, p.created_at
	          FROM users u
	          JOIN user_profile p ON u.id = p.user_id
	          WHERE u.id = $1 and role != 'c-admin' and u.deleted_at = 0`

	row := p.db.QueryRow(query1, userID)
	err = row.Scan(&user.Id, &user.Email, &user.Phone, &user.FirstName, &user.LastName, &user.Username, &user.Nationality,
		&user.Bio, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}
