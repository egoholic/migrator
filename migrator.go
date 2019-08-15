package migrator

type Migrator struct {
	Source string
	Target string
}

func New() *Migrator {

}

func (m *Migrator) Migrate()
