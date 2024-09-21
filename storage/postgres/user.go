package postgres

import (
	"auth-service/pkg/hashing"
	"auth-service/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"

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

	query := `INSERT INTO users (id, phone, email, password, first_name, last_name, username, nationality, bio) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	_, err := p.db.Exec(query, userID, req.Phone, req.Email, req.Password, req.FirstName, req.LastName, req.Username, req.Nationality, req.Bio)
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{Id: userID, Email: req.Email, FirstName: req.FirstName, LastName: req.LastName}, nil
}

func (p *UserRepo) GetProfile(req *pb.Id) (*pb.GetProfileResponse, error) {
	query := `SELECT id, email, phone, first_name, last_name, username, nationality, bio, 
	                 (SELECT COUNT(*) FROM follows WHERE following_id = users.id) AS follower_count, 
	                 (SELECT COUNT(*) FROM follows WHERE follower_id = users.id) AS following_count
	          FROM users
	          WHERE id = $1 AND role != 'admin' AND deleted_at = 0`

	row := p.db.QueryRow(query, req.UserId)
	var res pb.GetProfileResponse
	err := row.Scan(&res.UserId, &res.Email, &res.PhoneNumber, &res.FirstName, &res.LastName, &res.Username, &res.Nationality,
		&res.Bio, &res.FollowersCount, &res.FollowingCount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &res, nil
}

func (p *UserRepo) UpdateProfile(req *pb.UpdateProfileRequest) (*pb.UserResponse, error) {
	query := `UPDATE users 
	          SET first_name = $1, last_name = $2, username = $3, nationality = $4, bio = $5, profile_image = $6, updated_at = now()
	          WHERE id = $7 RETURNING id`

	_, err := p.db.Exec(query, req.FirstName, req.LastName, req.Username, req.Nationality, req.Bio, req.ProfileImage, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{Id: req.UserId, FirstName: req.FirstName, LastName: req.LastName, Email: ""}, nil
}

func (p *UserRepo) ChangePassword(req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	query := `SELECT password FROM users WHERE id = $1 AND deleted_at = 0`

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
	query := `UPDATE users SET profile_image = $1, updated_at = now() WHERE id = $2`
	_, err := p.db.Exec(query, req.Url, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.Void{}, nil
}

func (p *UserRepo) FetchUsers(req *pb.Filter) (*pb.UserResponses, error) {
	query := `SELECT id, email, first_name, last_name, username, created_at
	          FROM users
	          WHERE username ILIKE $1 AND role = 'user'
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

func (p *UserRepo) ListOfFollowing(req *pb.Id) (*pb.Follows, error) {
	followings := &pb.Follows{}

	query := `
		SELECT u.username, u.id
		FROM follows f
		JOIN users u ON f.following_id = u.id
		WHERE f.follower_id = $1;
    `

	err := p.db.Select(&followings, query, req.UserId)
	if err != nil {
		return nil, err
	}

	return followings, nil
}

func (p *UserRepo) ListOfFollowers(req *pb.Id) (*pb.Follows, error) {
	followers := &pb.Follows{}
	query := `
		SELECT u.username, u.id
		FROM follows f
		JOIN users u ON f.follower_id = u.id
		WHERE f.following_id = $1;
    `

	err := p.db.Select(&followers, query, req.UserId)
	if err != nil {
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
	query := `INSERT INTO follows (follower_id, following_id, followed_at)
	          VALUES ($1, $2, NOW())
	          RETURNING follower_id, following_id, followed_at`

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
		if errors.Is(err, sql.ErrNoRows) {
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

	query := `SELECT follower_id FROM follows
	          GROUP BY follower_id
	          ORDER BY COUNT(*) DESC LIMIT 1`

	var user pb.UserResponse
	err := p.db.QueryRow(query).Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get most popular user: %w", err)
	}

	query1 := `SELECT id, email, phone, first_name, last_name, username, nationality, bio, created_at
	          FROM users
	          WHERE id = $1 AND role != 'c-admin' AND deleted_at = 0`

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
