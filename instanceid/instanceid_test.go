package instanceid

import (
	"reflect"
	"testing"

	log "github.com/sirupsen/logrus"
)

func Test_splitOnPlus(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"one plain arg",
			args{" a"},
			[]string{" a"},
		},
		{
			"one plain arg with spaces",
			args{" a"},
			[]string{" a"},
		},
		{
			"one plain arg longer name",
			args{"someId"},
			[]string{"someId"},
		},
		{
			"two plain args",
			args{"a+b"},
			[]string{"a", "b"},
		},
		{
			"two longer plain args",
			args{"one+two"},
			[]string{"one", "two"},
		},

		{
			"three plain args",
			args{"a+b+c"},
			[]string{"a", "b", "c"},
		},
		{
			"three longer args",
			args{"one+two+three"},
			[]string{"one", "two", "three"},
		},
		{
			"three plain args with ()-1",
			args{"a+b+c()"},
			[]string{"a", "b", "c()"},
		},
		{
			"4 plain args with",
			args{"a+b+c+d"},
			[]string{"a", "b", "c", "d"},
		},
		{
			"three plain args with ()2",
			args{"a+(b()+c)"},
			[]string{"a", "(b()+c)"},
		},
		{
			"three structured args with ()",
			args{"a+(b()+c+ff(xx+zz))"},
			[]string{"a", "(b()+c+ff(xx+zz))"},
		},
		{
			"nested sum",
			args{"(a+b)"},
			[]string{"(a+b)"},
		},
		{
			"2 nested sum",
			args{"(x+y)+(a+b)"},
			[]string{"(x+y)", "(a+b)"},
		},
		{
			"invalid +",
			args{"+a"},
			[]string{"a"},
		},
		{
			"invalid postfix+",
			args{"a+"},
			[]string{"a"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitOnPlus(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitOnPlus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseMIID(t *testing.T) {
	type args struct {
		id string
	}

	log.SetLevel(log.TraceLevel)
	tests := []struct {
		name     string
		args     args
		wantMiid Miid
	}{
		{
			"simple",
			args{"msA/1.17/dev-123ab%3333s"},
			Miid{
				Sn: "msA",
				Vn: "1.17",
				Va: "dev-123ab",
				T:  3333,
			},
		},
		{
			"simple-short",
			args{"msA/1.17%3333s"},
			Miid{
				Sn: "msA",
				Vn: "1.17",
				T:  3333,
			},
		},
		{
			"simple-notSecond",
			args{"msA/1.17%3333"},
			Miid{
				Sn: "msA",
				Vn: "1.17",
				T:  3333,
			},
		},
		{
			"simple-notSecondNumber",
			args{"msA/1.17%333a"},
			Miid{
				Sn: "msA",
				Vn: "1.17",
				T:  0,
			},
		},
		{
			"toomanydelimiters",
			args{"msA/1.17/addInfo/surplusInfo%333s"},
			Miid{
				Sn: "msA",
				Vn: "1.17",
				T:  0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCiid := parseMIID(tt.args.id); !reflect.DeepEqual(gotCiid, tt.wantMiid) {
				t.Errorf("parseMIID() = %v, want %v", gotCiid, tt.wantMiid)
			}
		})
	}
}

func Test_parseCiid(t *testing.T) {
	log.SetLevel(log.TraceLevel)

	A := Miid{
		Sn: "A",
	}
	B := Miid{
		Sn: "B",
	}
	C := Miid{
		Sn: "C",
	}
	D := Miid{
		Sn: "D",
	}
	E := Miid{
		Sn: "E",
	}

	type args struct {
		id string
	}
	tests := []struct {
		name     string
		args     args
		wantCiid Ciid
	}{
		{
			"simple",
			args{
				"msA/1.17/dev-123ab%3333s",
			},
			Ciid{
				Miid: Miid{
					Sn: "msA",
					Vn: "1.17",
					Va: "dev-123ab",
					T:  3333,
				},
			},
		},
		{
			"one Call",
			args{
				"msA/1.17/dev-123ab%3333s(A)",
			},
			Ciid{
				Miid: Miid{Sn: "msA", Vn: "1.17", Va: "dev-123ab", T: 3333},
				Ciids: []Ciid{
					{
						Miid: Miid{Sn: "A"},
					},
				},
			},
		},
		{
			"one Call Plus another one",
			args{
				"msA/1.17/dev-123ab%3333s(A+B)",
			},
			Ciid{
				Miid: Miid{Sn: "msA", Vn: "1.17", Va: "dev-123ab", T: 3333},
				Ciids: []Ciid{
					{
						Miid: Miid{Sn: "A"},
					},
					{
						Miid: Miid{Sn: "B"},
					},
				},
			},
		},
		{
			"one Call Plus another one and one call",
			args{
				"msA/1.17/dev-123ab%3333s(A+B(C))",
			},
			Ciid{
				Miid: Miid{Sn: "msA", Vn: "1.17", Va: "dev-123ab", T: 3333},
				Ciids: []Ciid{
					{
						Miid: Miid{Sn: "A"},
					},
					{
						Miid: Miid{Sn: "B"},
						Ciids: []Ciid{
							{
								Miid: Miid{Sn: "C"},
							},
						},
					},
				},
			},
		},
		{
			"simple",
			args{
				"A(B+C(D+E))",
			},
			Ciid{A, []Ciid{
				{Miid: B},
				{C, []Ciid{
					{Miid: D},
					{Miid: E},
				}}}},
		},
		{
			"simple",
			args{
				"A(B)",
			},
			Ciid{A,
				[]Ciid{
					{Miid: B},
				},
			},
		},
		{
			"simple",
			args{
				"A(B+C)",
			},
			Ciid{A,
				[]Ciid{
					{Miid: B},
					{Miid: C},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCiid := parseCiid(tt.args.id); !reflect.DeepEqual(gotCiid, tt.wantCiid) {
				t.Errorf("parseCiid() = %v, want %v", gotCiid, tt.wantCiid)
			}
		})
	}
}
func Test_seperateFNameFromArg(t *testing.T) {
	type args struct {
		signature string
	}
	tests := []struct {
		name     string
		args     args
		wantName string
		wantArg  string
	}{
		{
			"simple",
			args{"A(B)"},
			"A",
			"B",
		},
		{
			"no Arg",
			args{"A"},
			"A", "",
		},
		{
			"empty Parenthesis",
			args{"A()"},
			"A",
			"",
		},
		{
			"no name",
			args{"(B)"},
			"",
			"B",
		},
		{
			"more complex",
			args{"A(B+C)"},
			"A",
			"B+C",
		},
		{
			"more complex, neste functions",
			args{"A(B(D)+C(E(F)))"},
			"A",
			"B(D)+C(E(F))",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotArg := seperateFNameFromArg(tt.args.signature)
			if gotName != tt.wantName {
				t.Errorf("seperateFNameFromArg() gotName = %v, want %v", gotName, tt.wantName)
			}
			if gotArg != tt.wantArg {
				t.Errorf("seperateFNameFromArg() gotArg = %v, want %v", gotArg, tt.wantArg)
			}
		})
	}
}

func Test_printCiid(t *testing.T) {
	type args struct {
		ciid Ciid
	}
	log.SetLevel(log.TraceLevel)
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"simple",
			args{parseCiid("A(B+C)")},
			49,
		},
		{
			"simple2",
			args{parseCiid("A(B+C(D))")},
			70,
		},
		{
			"iid1",
			args{parseCiid("msA/1.1/abs%22s(msB/2.0/xxxx%333s+C(D))")},
			95,
		},
		{
			"iid2",
			args{parseCiid("msA/1.1/abs%22s(msB/2.0/xxxx%333s+msC/1.1%22s(D))")},
			107,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(PrintCiid(tt.args.ciid)); got != tt.want {
				t.Errorf("printCiid() = %v, want %v", got, tt.want)
				t.Errorf("theTree() = \n%v", PrintCiid(tt.args.ciid))
			}
		})
	}
}

func TestCiid_String(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"simpleMiid",
			"msA/1.1%22s",
		},
		{
			"fullMiid",
			"msA/1.1/feature-branch-22aabbcc%22s",
		},
		{
			"emptyMiid",
			"",
		},
		{
			"justSimpleMidd",
			"A",
		},
		{
			"fullMiidOneCiid",
			"msA/1.1/feature-branch-22aabbcc%22s(msB)",
		},
		{
			"fullMiidTwoCiid",
			"msA/1.1/feature-branch-22aabbcc%22s(msB+msC)",
		},
		{
			"complexFunc",
			"A(B(C+D)+D(E))",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewCiid(tt.want)
			if got := mock.String(); got != tt.want {
				t.Errorf("Ciid.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCiid(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name     string
		args     args
		wantCiid Ciid
	}{
		{
			"simpleMiid",
			args{"msA/1.1%22"},
			Ciid{
				Miid{
					Sn: "msA",
					Vn: "1.1",
					Va: "",
					T:  22,
				},
				nil,
			},
		},
		{
			"fullMiid",
			args{"msA/1.1/feature-branch-22aabbcc%22"},
			Ciid{
				Miid: Miid{
					Sn: "msA",
					Vn: "1.1",
					Va: "feature-branch-22aabbcc",
					T:  22,
				},
				Ciids: nil,
			},
		},
		{
			"emptyMiid",
			args{""},
			Ciid{
				Miid: Miid{
					Sn: "",
					Vn: "",
					Va: "",
					T:  0,
				},
				Ciids: nil,
			},
		},
		{
			"fullMiidOneCiid",
			args{"msA/1.1/feature-branch-22aabbcc%22(msB)"},
			Ciid{
				Miid: Miid{
					Sn: "msA",
					Vn: "1.1",
					Va: "feature-branch-22aabbcc",
					T:  22,
				},
				Ciids: []Ciid{
					{
						Miid: Miid{
							Sn: "msB",
							Vn: "",
							Va: "",
							T:  0,
						},
						Ciids: nil,
					},
				},
			},
		},
		{
			"fullMiidTwoCiid",
			args{"msA/1.1/feature-branch-22aabbcc%22s(msB+msC)"},
			Ciid{
				Miid: Miid{
					Sn: "msA",
					Vn: "1.1",
					Va: "feature-branch-22aabbcc",
					T:  22,
				},
				Ciids: []Ciid{
					{
						Miid: Miid{
							Sn: "msB",
							Vn: "",
							Va: "",
							T:  0,
						},
						Ciids: nil,
					},
					{
						Miid: Miid{
							Sn: "msC",
							Vn: "",
							Va: "",
							T:  0,
						},
						Ciids: nil,
					},
				},
			},
		},
		{
			"complexFunc",
			args{"A(B(C+D)+D(E)"},
			Ciid{
				Miid: Miid{
					Sn: "A",
				},
				Ciids: []Ciid{
					{
						Miid: Miid{
							Sn: "B",
							Vn: "",
							Va: "",
							T:  0,
						},
						Ciids: []Ciid{
							{
								Miid: Miid{
									Sn: "C",
									Vn: "",
									Va: "",
									T:  0,
								},
								Ciids: nil,
							},
							{
								Miid: Miid{
									Sn: "D",
									Vn: "",
									Va: "",
									T:  0,
								},
								Ciids: nil,
							},
						},
					},
					{
						Miid: Miid{
							Sn: "D",
							Vn: "",
							Va: "",
							T:  0,
						},
						Ciids: []Ciid{
							{
								Miid: Miid{
									Sn: "E",
									Vn: "",
									Va: "",
									T:  0,
								},
								Ciids: nil,
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCiid := NewCiid(tt.args.id); !reflect.DeepEqual(gotCiid, tt.wantCiid) {
				t.Errorf("NewCiid() = %#v, want %#v", gotCiid, tt.wantCiid)
			}
		})
	}
}

func TestNewMiid(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name     string
		args     args
		wantMiid Miid
	}{
		{
			"CiidNotExpected",
			args{"msA/1.1/feature-branch-22aabbcc%22s(msB+msC)"},
			Miid{
				Sn: "msA",
				Vn: "1.1",
				Va: "feature-branch-22aabbcc",
				T:  22,
			},
		},
		{
			"simple",
			args{"msA"},
			Miid{
				Sn: "msA",
				Vn: "",
				Va: "",
				T:  0,
			},
		},
		{
			"complex",
			args{"msA/1.1/asdfasdf-asdfasdf%22s"},
			Miid{
				Sn: "msA",
				Vn: "1.1",
				Va: "asdfasdf-asdfasdf",
				T:  22,
			},
		},
		{
			"no clue",
			args{"This is some text"},
			Miid{
				Sn: "This is some text",
				Vn: "",
				Va: "",
				T:  0,
			},
		},
		{
			"no clue",
			args{"(/)"},
			Miid{
				Sn: "(",
				Vn: ")",
				Va: "",
				T:  0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotMiid := NewMiid(tt.args.id); !reflect.DeepEqual(gotMiid, tt.wantMiid) {
				t.Errorf("NewMiid() = %v, want %v", gotMiid, tt.wantMiid)
			}
		})
	}
}
