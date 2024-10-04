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
	"strconv"
	"strings"

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

	query := `INSERT INTO users (id, phone, email, password, first_name, last_name, username, country, bio, role) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := p.db.Exec(query, userID, req.Phone, req.Email, req.Password,
		req.FirstName, req.LastName, req.Username, req.Nationality, req.Bio, req.Role)
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{Id: userID, Email: req.Email, FirstName: req.FirstName, LastName: req.LastName}, nil
}

func (p *UserRepo) GetProfile(req *pb.Id) (*pb.GetProfileResponse, error) {
	query := `SELECT id, email, phone, first_name, last_name, username, country, bio, 
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
	setValues := make([]string, 0, 6)
	args := make([]interface{}, 0, 7)
	count := 1
	if req.FirstName != "" {
		setValues = append(setValues, "first_name = $"+strconv.Itoa(count))
		count++
		args = append(args, req.FirstName)
	}
	if req.LastName != "" {
		setValues = append(setValues, "last_name = $"+strconv.Itoa(count))
		count++
		args = append(args, req.LastName)
	}
	if req.Username != "" {
		setValues = append(setValues, "username = $"+strconv.Itoa(count))
		count++
		args = append(args, req.Username)
	}
	if req.Nationality != "" {
		setValues = append(setValues, "country = $"+strconv.Itoa(count))
		count++
		args = append(args, req.Nationality)
	}
	if req.Bio != "" {
		setValues = append(setValues, "bio = $"+strconv.Itoa(count))
		count++
		args = append(args, req.Bio)
	}

	setValuesStr := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE users 
	                       SET %s, updated_at = now()
	                       WHERE id = $%d RETURNING id`, setValuesStr, len(args)+1)

	args = append(args, req.UserId)

	_, err := p.db.Exec(query, args...)
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
	query := `UPDATE users SET profile_image = $1, updated_at = now() WHERE id = $2 and deleted_at = 0`
	_, err := p.db.Exec(query, req.Url, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.Void{}, nil
}

func (p *UserRepo) FetchUsers(req *pb.Filter) (*pb.UserResponses, error) {
	where := "WHERE role = 'user'"
	if req.FirstName != "" {
		where += fmt.Sprintf(" AND username ILIKE '%s%%'", req.FirstName)
	}

	query := fmt.Sprintf(`SELECT id, email, first_name, last_name, username, created_at
	          FROM users
	          %s
	          LIMIT $1 OFFSET $2`, where)

	rows, err := p.db.Query(query, req.Limit, (req.Page-1)*req.Limit)
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

	err := p.db.Select(&followings.Following, query, req.UserId)
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

	err := p.db.Select(&followers.Following, query, req.UserId)
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

	query1 := `SELECT id, email, phone, first_name, last_name, username, country, bio, created_at
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

// AddNationality adds a new nationality to the database
func (p *UserRepo) AddNationality(req *pb.Nationality) (*pb.Void, error) {

	query := `INSERT INTO nationality (description, name) VALUES ($1, $2)`
	_, err := p.db.Exec(query, req.Description, req.Name)
	if err != nil {
		return nil, err
	}

	return &pb.Void{}, nil
}

// GetNationalityById retrieves a nationality by its ID
func (p *UserRepo) GetNationalityById(req *pb.NId) (*pb.Nationality, error) {
	query := `SELECT id,description, name FROM nationality WHERE id = $1`

	var nationality pb.Nationality
	err := p.db.QueryRow(query, req.Id).Scan(&nationality.Id, &nationality.Description, &nationality.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("nationality not found")
		}
		return nil, err
	}

	return &nationality, nil
}

// ListNationalities retrieves a list of nationality with pagination and filter
func (p *UserRepo) ListNationalities(req *pb.Pagination) (*pb.Nationalities, error) {
	where := "WHERE TRUE"
	if req.Name != "" {
		where += fmt.Sprintf(" AND LOWER(name) LIKE '%%%s%%'", strings.ToLower(req.Name))
	}

	if req.Description != "" {
		where += fmt.Sprintf(" AND LOWER(description) LIKE '%%%s%%'", strings.ToLower(req.Description))
	}

	query := fmt.Sprintf(`SELECT id, name, description FROM nationality %s LIMIT $1 OFFSET $2`, where)
	rows, err := p.db.Query(query, req.Limit, (req.Page-1)*req.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nationalities []*pb.Nationality
	for rows.Next() {
		var nationality pb.Nationality
		if err := rows.Scan(&nationality.Id, &nationality.Name, &nationality.Description); err != nil {
			return nil, err
		}
		nationalities = append(nationalities, &nationality)
	}

	return &pb.Nationalities{Nationalities: nationalities}, nil
}

// UpdateNationality updates the name of an existing nationality
func (p *UserRepo) UpdateNationality(req *pb.Nationality) (*pb.Void, error) {
	setValues := make([]string, 0, 2)
	args := make([]interface{}, 0, 3)
	count := 1
	if req.Name != "" {
		setValues = append(setValues, "name = $"+strconv.Itoa(count))
		args = append(args, req.Name)
		count++
	}
	if req.Description != "" {
		setValues = append(setValues, "description = $"+strconv.Itoa(count))
		args = append(args, req.Description)
		count++
	}

	setValuesStr := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE nationality
	                       SET %s
	                       WHERE id = $%d`, setValuesStr, count)

	args = append(args, req.Id)

	_, err := p.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return &pb.Void{}, nil
}

// DeleteNationality deletes a nationality from the database
func (p *UserRepo) DeleteNationality(req *pb.NId) (*pb.Void, error) {
	query := `DELETE FROM nationality WHERE id = $1`
	_, err := p.db.Exec(query, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.Void{}, nil
}
