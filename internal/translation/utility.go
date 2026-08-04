package translation

import (
	"time"

	repo "github.com/Will-Gould/psmp-api/internal/adapters/mysql/sqlc"
)

func LedgerTranslateToBlocks(actions []repo.Action) []Block {
	blocks := []Block{}
	for _, a := range actions {
		blocks = append(blocks, Block{
			Time:   a.Time.UnixMilli(),
			User:   a.PlayerID.Int32,
			World:  a.WorldID,
			X:      a.X,
			Y:      a.Y,
			Z:      a.Z,
			Object: a.ObjectID,
			Action: a.ActionID,
		})
	}
	return blocks
}

func GriefLoggerTranslateToBlocks(glBlocks []repo.Block) []Block {
	blocks := []Block{}
	for _, b := range glBlocks {
		blocks = append(blocks, Block{
			Time:   b.Time,
			User:   b.User,
			World:  b.Level,
			X:      b.X,
			Y:      b.Y,
			Z:      b.Z,
			Object: b.Type,
			Action: b.Action,
		})
	}
	return blocks
}

func GriefLoggerTranslateSessions(glSessions []repo.Session) []Session {
	sessions := []Session{}
	for _, s := range glSessions {
		sessions = append(sessions, Session{
			ID:     s.User,
			Time:   s.Time,
			Action: s.Action,
		})
	}
	return sessions
}

func PsmpstatsTranslateSessions(pSessions []repo.PsmpstatsSession) []Session {
	sessions := []Session{}
	for _, s := range pSessions {
		sessions = append(sessions, Session{
			ID:     s.PlayerID,
			Time:   time.Unix(int64(s.Time), 0).UnixMilli(),
			Action: s.Action,
		})
	}
	return sessions
}
