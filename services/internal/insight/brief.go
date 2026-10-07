package insight

import (
	"fmt"
	"strings"
)

// Brief is what the writer is told about one moment of a match. Every fact in
// it comes from the Flow engine and the outcome model: the text only retells
// them, so it cannot invent a shot that never happened.
type Brief struct {
	Key         string // the moment, the same in every replay: "3754314:m30"
	Minute      int
	Competition string
	Home, Away  Team
	Events      []string // the key events so far, oldest first: "23′ гол — Arsenal (Mesut Özil)"
	Chances     Chances  // the outcome model now
	PreMatch    Chances  // and before kick-off
}

// Team is one side at that moment.
type Team struct {
	Name    string
	Goals   int
	Flow    int      // 0..100
	Delta10 int      // the change of Flow over the last 10 minutes
	Factors []string // what its Flow is made of: "удары +14 (4 за 7 мин)"
}

// Chances are whole percents of a home win, a draw and an away win.
type Chances struct {
	Home, Draw, Away int
}

// The instructions are the same for every moment; the facts change. A writer
// sends them apart, as the system instruction and the message.
const instructions = `Ты — футбольный аналитик приложения FlowScore. Тебе дают факты об одном моменте матча, ты пишешь короткий разбор для болельщика.

Flow — метрика давления команды от 0 до 100: она растёт от ударов, ключевых передач, угловых и владения и затухает со временем. Факторы показывают, из чего сейчас складывается Flow команды. Шансы — прогноз модели исхода.

Правила:
- Только факты из сообщения. Не придумывай игроков, моменты, травмы, составы и статистику, которых нет, и не делай выводов о том, как забиты голы.
- Заголовок: до 6 слов, без точки в конце.
- Текст: 2–3 коротких предложения. Что происходит сейчас, почему (опирайся на факторы и события) и что это значит для исхода.
- Числа бери из фактов, не пересчитывай. Не повторяй все числа подряд: одного-двух хватит.
- Тон спокойный, как у аналитика в эфире. Без восклицательных знаков и эмодзи.`

// The answer is JSON of this shape: the API holds the model to it
var textSchema = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"title": map[string]any{"type": "string"},
		"text":  map[string]any{"type": "string"},
	},
	"required":             []string{"title", "text"},
	"additionalProperties": false,
}

// prompt is the message: the facts as plain lines, in Russian like the text
// it asks for.
func (b Brief) prompt() string {
	var s strings.Builder
	fmt.Fprintf(&s, "Турнир: %s\n", b.Competition)
	fmt.Fprintf(&s, "Минута: %d\n", b.Minute)
	fmt.Fprintf(&s, "Счёт: %s %d:%d %s\n", b.Home.Name, b.Home.Goals, b.Away.Goals, b.Away.Name)

	for _, t := range []Team{b.Home, b.Away} {
		fmt.Fprintf(&s, "\nFlow %s: %d (за 10 минут %+d)\n", t.Name, t.Flow, t.Delta10)
		for _, f := range t.Factors {
			fmt.Fprintf(&s, "- %s\n", f)
		}
	}

	s.WriteString("\nКлючевые события:\n")
	if len(b.Events) == 0 {
		s.WriteString("- пока нет\n")
	}
	for _, e := range b.Events {
		fmt.Fprintf(&s, "- %s\n", e)
	}

	fmt.Fprintf(&s, "\nШансы сейчас: П1 %d%%, X %d%%, П2 %d%%\n", b.Chances.Home, b.Chances.Draw, b.Chances.Away)
	fmt.Fprintf(&s, "Шансы до матча: П1 %d%%, X %d%%, П2 %d%%\n", b.PreMatch.Home, b.PreMatch.Draw, b.PreMatch.Away)
	return s.String()
}
