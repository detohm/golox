package golox

import (
	"testing"
)

func TestAstPrinter_Print(t *testing.T) {
	expression := NewBinary(
		NewUnary(
			NewToken(TkMinus, "-", nil, 1),
			NewLiteral(123),
		),
		NewToken(TkStar, "*", nil, 1),
		NewGrouping(
			NewLiteral(45.67),
		),
	)

	t.Run("Verify printer result", func(t *testing.T) {
		result := NewAstPrinter().Print(expression)
		expected := "(* (- 123) (group 45.67))"
		if result != expected {
			t.Errorf("AstPrinter_Print result %s, expected %s",
				result,
				expected,
			)
		}
	})
}
