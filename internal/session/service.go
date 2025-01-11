package session

import (
	"errors"

	"github.com/CaribouBlue/top-spot/internal/music"
	"github.com/google/uuid"
)

var (
	ErrSubmissionNotFound = errors.New("submission not found")
	ErrVoteNotFound       = errors.New("vote not found")
	ErrPlaylistNotFound   = errors.New("playlist not found")
	ErrPlaylistExists     = errors.New("playlist already exists")
)

type SessionService interface {
	GetAll() ([]*Session, error)
	GetOne(sessionId int64) (*Session, error)
	Create(*Session) error

	GetSubmissions(sessionId int64) ([]Submission, error)
	GetSubmission(sessionId int64, submissionId string) (*Submission, error)
	AddSubmission(sessionId int64, submission *Submission) error
	RemoveSubmission(sessionId int64, submissionId string) error

	GetVotes(sessionId int64) ([]Vote, error)
	GetVote(sessionId int64, voteId string) (*Vote, error)
	AddVote(sessionId int64, vote *Vote) error
	RemoveVote(sessionId int64, voteId string) error

	GetPlaylist(sessionId int64, userId int64) (*Playlist, error)
	AddPlaylist(sessionId int64, playlistId string, userId int64) error
}

type sessionService struct {
	store          SessionStore
	musicService   music.MusicService
	collectionName string
}

func NewSessionService(store SessionStore, musicService music.MusicService) SessionService {
	return &sessionService{
		store:          store,
		musicService:   musicService,
		collectionName: "session",
	}
}

func (s *sessionService) GetAll() ([]*Session, error) {
	return s.store.GetSessions()
}

func (s *sessionService) GetOne(sessionId int64) (*Session, error) {
	return s.store.GetSession(sessionId)
}

func (s *sessionService) Create(session *Session) error {
	return s.store.CreateSession(session)
}

func (s *sessionService) GetSubmissions(sessionId int64) ([]Submission, error) {
	session, err := s.store.GetSession(sessionId)
	if err != nil {
		return nil, err
	}

	return session.Submissions, nil
}

func (s *sessionService) GetSubmission(sessionId int64, submissionId string) (*Submission, error) {
	session, err := s.store.GetSession(sessionId)
	if err != nil {
		return nil, err
	}

	for _, submission := range session.Submissions {
		if submission.Id == submissionId {
			return &submission, nil
		}
	}

	return nil, ErrSubmissionNotFound
}

func (s *sessionService) AddSubmission(sessionId int64, submission *Submission) error {
	session, err := s.store.GetSession(sessionId)
	if err != nil {
		return err
	}

	submission.Id = uuid.New().String()

	session.Submissions = append(session.Submissions, *submission)

	return s.store.UpdateSession(session)
}

func (s *sessionService) RemoveSubmission(sessionId int64, submissionId string) error {
	session, err := s.store.GetSession(sessionId)
	if err != nil {
		return err
	}

	for i, submission := range session.Submissions {
		if submission.Id == submissionId {
			session.Submissions = append(session.Submissions[:i], session.Submissions[i+1:]...)
			return s.store.UpdateSession(session)
		}
	}

	return ErrSubmissionNotFound
}

func (s *sessionService) GetVotes(sessionId int64) ([]Vote, error) {
	session, err := s.store.GetSession(sessionId)
	if err != nil {
		return nil, err
	}

	return session.Votes, nil
}

func (s *sessionService) GetVote(sessionId int64, voteId string) (*Vote, error) {
	session, err := s.store.GetSession(sessionId)
	if err != nil {
		return nil, err
	}

	for _, vote := range session.Votes {
		if vote.Id == voteId {
			return &vote, nil
		}
	}

	return nil, ErrVoteNotFound
}

func (s *sessionService) AddVote(sessionId int64, vote *Vote) error {
	session, err := s.store.GetSession(sessionId)
	if err != nil {
		return err
	}

	session.Votes = append(session.Votes, *vote)

	return s.store.UpdateSession(session)
}

func (s *sessionService) RemoveVote(sessionId int64, voteId string) error {
	session, err := s.store.GetSession(sessionId)
	if err != nil {
		return err
	}

	for i, vote := range session.Votes {
		if vote.Id == voteId {
			session.Votes = append(session.Votes[:i], session.Votes[i+1:]...)
			return s.store.UpdateSession(session)
		}
	}

	return ErrVoteNotFound
}

// GetPlaylist(sessionId int64, userId int64) (Playlist, error)
// AddPlaylist(sessionId int64, playlistId string, userId int64) error

func (s *sessionService) GetPlaylist(sessionId int64, userId int64) (*Playlist, error) {
	session, err := s.store.GetSession(sessionId)
	if err != nil {
		return nil, err
	}

	for _, playlist := range session.Playlists {
		if playlist.UserId == userId {
			return &playlist, nil
		}
	}

	return nil, ErrPlaylistNotFound
}

func (s *sessionService) AddPlaylist(sessionId int64, playlistId string, userId int64) error {
	session, err := s.store.GetSession(sessionId)
	if err != nil {
		return err
	}

	for _, playlist := range session.Playlists {
		if playlist.UserId == userId {
			return ErrPlaylistExists
		}
	}

	session.Playlists = append(session.Playlists, Playlist{
		Id:     playlistId,
		UserId: userId,
	})

	return s.store.UpdateSession(session)
}
