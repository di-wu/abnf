package operators

import (
	"fmt"
	"testing"
)

func (n *Node) equals(other *Node) error {
	if n == nil || other == nil {
		return fmt.Errorf("one of the nodes is nil")
	}

	if n.Key != other.Key {
		return fmt.Errorf("keys do not match: %s, %s", n.Key, other.Key)
	}
	if len(n.Value) != len(other.Value) {
		return fmt.Errorf("values do not match: %s, %s", string(n.Value), string(other.Value))
	}
	if len(n.Children) != len(other.Children) {
		return fmt.Errorf("not an equal amount of children: %d, %d", len(n.Children), len(other.Children))
	}

	for i := range n.Value {
		if n.Value[i] != other.Value[i] {
			return fmt.Errorf("value does not match, index %d", i)
		}
	}

	for i := range n.Children {
		if err := n.Children[i].equals(other.Children[i]); err != nil {
			return err
		}
	}

	return nil
}

func (as Alternatives) equals(other Alternatives) error {
	if len(as) != len(other) {
		return fmt.Errorf("not an equal amount of nodes: %d, %d", len(as), len(other))
	}

	for i, alternative := range as {
		if err := alternative.equals(other[i]); err != nil {
			return err
		}
	}
	return nil
}

func TestNode(t *testing.T) {
	for _, test := range []struct {
		name    string
		rule    Operator
		str     string
		correct Alternatives
	}{
		{
			name: "Rune",
			rule: a,
			str:  "a",
			correct: Alternatives{
				{
					Key:   "a",
					Value: []rune("a"),
				},
			},
		},
		{
			name: "RuneEmpty",
			rule: a,
			str:  "b",
		},
		{
			name: "String",
			rule: String(`s`, "str"),
			str:  "Str",
			correct: Alternatives{
				{
					Key:   "s",
					Value: []rune("Str"),
				},
			},
		},
		{
			name: "StringEmpty",
			rule: String(`s`, "str"),
			str:  "rts",
		},
		{
			name: "StringCS",
			rule: StringCS(`cs`, "str"),
			str:  "str",
			correct: Alternatives{
				{
					Key:   "cs",
					Value: []rune("str"),
				},
			},
		},
		{
			name: "StringCSEmpty",
			rule: StringCS(`cs`, "str"),
			str:  "rts",
		},
		{
			name: "Range",
			rule: Range(`r`, 'a', 'z'),
			str:  "x",
			correct: Alternatives{
				{
					Key:   "r",
					Value: []rune("x"),
				},
			},
		},
		{
			name: "RangeEmpty",
			rule: Range(`r`, 'a', 'z'),
			str:  "0",
		},
		{
			name: "Optional",
			rule: Optional(`o`, a),
			str:  "a",
			correct: Alternatives{
				{
					Key:   "o",
					Value: []rune("a"),
					Children: Children{
						{
							Key:   "a",
							Value: []rune("a"),
						},
					},
				},
				{
					Key:   "o",
					Value: []rune(""),
				},
			},
		},
		{
			name: "OptionalEmpty",
			rule: Optional(`o`, a),
			correct: Alternatives{
				{
					Key:   "o",
					Value: []rune(""),
				},
			},
		},
		{
			name: "Repeat0",
			rule: RepeatN(`r0`, 0, a),
			correct: Alternatives{
				{
					Key:   "r0",
					Value: []rune(""),
				},
			},
		},
		{
			name: "Repeat1",
			rule: RepeatN(`r1`, 1, a),
			str:  "aaa",
			correct: Alternatives{
				{
					Key:   "r1",
					Value: []rune("a"),
				},
			},
		},
		{
			name: "Repeat1Empty",
			rule: RepeatN(`r1`, 1, a),
			str:  "bbb",
		},
		{
			name: "Repeat0Inf",
			rule: Repeat0Inf(`r0i`, a),
			str:  "aa",
			correct: Alternatives{
				{
					Key:   "r0i",
					Value: []rune("aa"),
					Children: Children{
						{
							Key:   "r0i",
							Value: []rune("a"),
							Children: Children{
								{
									Key:   "r0i",
									Value: []rune(""),
								},
							},
						},
						{
							Key:   "r0i",
							Value: []rune(""),
						},
					},
				},
				{
					Key:   "r0i",
					Value: []rune("a"),
					Children: Children{
						{
							Key:   "r0i",
							Value: []rune(""),
						},
					},
				},
				{
					Key:   "r0i",
					Value: []rune(""),
				},
			},
		},
		{
			name: "Repeat0InfEmpty",
			rule: Repeat0Inf(`r0i`, a),
			str:  "",
			correct: Alternatives{
				{
					Key:   "r0i",
					Value: []rune(""),
				},
			},
		},
		{
			name: "Repeat1Inf",
			rule: Repeat1Inf(`r1i`, a),
			str:  "aa",
			correct: Alternatives{
				{
					Key:   "r1i",
					Value: []rune("aa"),
					Children: Children{
						{
							Key:   "r1i",
							Value: []rune("a"),
						},
					},
				},
				{
					Key:   "r1i",
					Value: []rune("a"),
				},
			},
		},
		{
			name: "Repeat1InfEmpty",
			rule: Repeat1Inf(`r1i`, a),
			str:  "bbb",
		},
		{
			name: "RepeatOptional",
			rule: RepeatOptional(`ro`, a),
			str:  "a",
			correct: Alternatives{
				{
					Key:   "ro",
					Value: []rune("a"),
					Children: Children{
						{
							Key:   "ro",
							Value: []rune(""),
						},
					},
				},
				{
					Key:   "ro",
					Value: []rune(""),
				},
			},
		},
		{
			name: "RepeatOptionalEmpty",
			rule: RepeatOptional(`ro`, a),
			str:  "bbb",
			correct: Alternatives{
				{
					Key:   "ro",
					Value: []rune(""),
				},
			},
		},
		{
			name: "Concat",
			rule: Concat(`c`, a, b, c),
			str:  "abc",
			correct: Alternatives{
				{
					Key:   "c",
					Value: []rune("abc"),
					Children: Children{
						{
							Key:   "a",
							Value: []rune("a"),
						},
						{
							Key:   "b",
							Value: []rune("b"),
						},
						{
							Key:   "c",
							Value: []rune("c"),
						},
					},
				},
			},
		},
		{
			name: "ConcatEmpty",
			rule: Concat(`c`, a, b, c),
			str:  "cba",
		},
		{
			name: "ConcatRepeat",
			rule: Concat(`cr`,
				Repeat0Inf(`r0i1`, a),
				Repeat0Inf(`r0i2`, a),
			),
			str: "a",
			correct: Alternatives{
				{
					Key:   "cr",
					Value: []rune("a"),
					Children: Children{
						{
							Key:   "r0i1",
							Value: []rune("a"),
							Children: Children{
								{
									Key:   "r0i1",
									Value: []rune(""),
								},
							},
						},
						{
							Key:   "r0i2",
							Value: []rune(""),
						},
					},
				},
				{
					Key:   "cr",
					Value: []rune("a"),
					Children: Children{
						{
							Key:   "r0i1",
							Value: []rune(""),
						},
						{
							Key:   "r0i2",
							Value: []rune("a"),
							Children: Children{
								{
									Key:   "r0i2",
									Value: []rune(""),
								},
							},
						},
					},
				},
				{
					Key:   "cr",
					Value: []rune(""),
					Children: Children{
						{
							Key:   "r0i1",
							Value: []rune(""),
						},
						{
							Key:   "r0i2",
							Value: []rune(""),
						},
					},
				},
			},
		},
		{
			name: "AltsOption1",
			rule: Alts(`a`, a, b),
			str:  "a",
			correct: Alternatives{
				{
					Key:   "a",
					Value: []rune("a"),
					Children: Children{
						{
							Key:   "a",
							Value: []rune("a"),
						},
					},
				},
			},
		},
		{
			name: "AltsOption2",
			rule: Alts(`a`, a, b),
			str:  "b",
			correct: Alternatives{
				{
					Key:   "a",
					Value: []rune("b"),
					Children: Children{
						{
							Key:   "b",
							Value: []rune("b"),
						},
					},
				},
			},
		},
		{
			name: "AltsEmpty",
			rule: Alts(`a`, a, b),
			str:  "c",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			nodes := test.rule([]rune(test.str))
			if err := nodes.equals(test.correct); err != nil {
				for _, node := range nodes {
					fmt.Print(node.StringRecursive())
				}
				t.Error(err)
			}
		})
	}
}

func TestBestAlternative(t *testing.T) {
	t.Run(`Simple`, func(t *testing.T) {
		str := "aaa"
		nodes := Repeat0Inf(``, a)([]rune(str))
		if best := nodes.Best(); best.String() != str {
			t.Error("did not get best alternative")
		}
	})
}