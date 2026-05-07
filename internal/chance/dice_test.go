package chance

import (
	"fmt"
	"testing"
)

var dieTests = []struct {
	name  string
	sided int
	amt   int
}{
	{name: "one d6", sided: 6, amt: 1},
	{name: "two d6", sided: 6, amt: 2},
	{name: "three d6", sided: 6, amt: 3},
}

func TestDie_Roll(t *testing.T) {
	d := Die{Sided: 6}
	for i := 0; i < 100; i++ {
		got := d.Roll()
		if got < 1 || got > 6 {
			t.Errorf("Die.Roll() = %d; want in range [1, 6]", got)
		}
	}
}

func TestCup_Add(t *testing.T) {
	c := NewCup(Die{Sided: 6})
	c.Add(Die{Sided: 4}, Die{Sided: 8})
	got := c.Throw()
	if len(got.Rolls) != 3 {
		t.Errorf("len(Rolls) = %d; want 3", len(got.Rolls))
	}
}

func TestCup_Shake(t *testing.T) {
	// Shake just shuffles, hard to test without seeding or large samples,
	// but we can at least ensure it doesn't crash and preserves dice count.
	c := NewCup(Die{Sided: 6}, Die{Sided: 4})
	c.Shake()
	got := c.Throw()
	if len(got.Rolls) != 2 {
		t.Errorf("len(Rolls) = %d; want 2", len(got.Rolls))
	}
}

func TestCup_Throw(t *testing.T) {
	for _, tt := range dieTests {
		t.Run(tt.name, func(t *testing.T) {
			dice := make([]Die, tt.amt)
			for i := 0; i < tt.amt; i++ {
				dice[i] = Die{Sided: tt.sided}
			}
			c := NewCup(dice...)
			got := c.Throw()
			if len(got.Rolls) != tt.amt {
				t.Errorf("number of dice wanted = %v, got %v", tt.amt, len(got.Rolls))
			}
			var total int
			for _, val := range got.Rolls {
				if val < 1 || val > tt.sided {
					t.Errorf("value should be in range 1 to %d, got %d", tt.sided, val)
				}
				total += val
			}
			if total != got.Total {
				t.Errorf("total of rolls should be %d, got %d", total, got.Total)
			}
		})
	}
}

func BenchmarkCup_Throw(b *testing.B) {
	for _, d := range dieTests {
		dice := make([]Die, d.amt)
		for i := 0; i < d.amt; i++ {
			dice[i] = Die{Sided: d.sided}
		}
		c := NewCup(dice...)
		b.Run(fmt.Sprintf("%dd%d", d.amt, d.sided), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = c.Throw()
			}
		})
	}
}
