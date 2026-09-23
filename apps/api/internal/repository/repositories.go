package repository

import "github.com/chandankrr/loreline/internal/server"

type Repositories struct {
	User         *UserRepository
	Session      *SessionRepository
	Account      *AccountRepository
	Verification *VerificationRepository
	Folder       *FolderRepository
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{
		User:         NewUserRepository(s),
		Session:      NewSessionRepository(s),
		Account:      NewAccountRepository(s),
		Verification: NewVerificationRepository(s),
		Folder:       NewFolderRepository(s),
	}
}
