package parser_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/pboyd04/godata/filter/lexer"
	"github.com/pboyd04/godata/filter/parser"
)

type testData struct {
	input             string
	expectedOperation parser.Operation
}

//nolint:gochecknoglobals
var testCases = []testData{
	{
		input:             "true",
		expectedOperation: parser.Operation{Operator: parser.Operator(lexer.TokenTrue)},
	},
	{
		input:             "false",
		expectedOperation: parser.Operation{Operator: parser.Operator(lexer.TokenFalse)},
	},
	{
		input: "Name eq 'Milk'",
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Equals),
			Operands: []parser.Operand{
				lexer.Token{Text: "Name", Type: lexer.UnquotedString},
				lexer.Token{Text: "'Milk'", Type: lexer.SingleQuotedString},
			},
		},
	},
	{
		input: "(Name eq 'Milk')",
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Equals),
			Operands: []parser.Operand{
				lexer.Token{Text: "Name", Type: lexer.UnquotedString},
				lexer.Token{Text: "'Milk'", Type: lexer.SingleQuotedString},
			},
		},
	},
	{
		input: "Name ne 'Milk'",
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.NotEquals),
			Operands: []parser.Operand{
				lexer.Token{Text: "Name", Type: lexer.UnquotedString},
				lexer.Token{Text: "'Milk'", Type: lexer.SingleQuotedString},
			},
		},
	},
	{
		input: "Name gt 'Milk'",
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.GreaterThan),
			Operands: []parser.Operand{
				lexer.Token{Text: "Name", Type: lexer.UnquotedString},
				lexer.Token{Text: "'Milk'", Type: lexer.SingleQuotedString},
			},
		},
	},
	{
		input: "Name ge 'Milk'",
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.GreaterThanOrEqual),
			Operands: []parser.Operand{
				lexer.Token{Text: "Name", Type: lexer.UnquotedString},
				lexer.Token{Text: "'Milk'", Type: lexer.SingleQuotedString},
			},
		},
	},
	{
		input: "Name lt 'Milk'",
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.LessThan),
			Operands: []parser.Operand{
				lexer.Token{Text: "Name", Type: lexer.UnquotedString},
				lexer.Token{Text: "'Milk'", Type: lexer.SingleQuotedString},
			},
		},
	},
	{
		input: "Name le 'Milk'",
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.LessThanOrEqual),
			Operands: []parser.Operand{
				lexer.Token{Text: "Name", Type: lexer.UnquotedString},
				lexer.Token{Text: "'Milk'", Type: lexer.SingleQuotedString},
			},
		},
	},
	{
		input: "Name eq 'Milk' and Price lt 2.55",
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.And),
			Operands: []parser.Operand{
				&parser.Operation{
					Operator: parser.Operator(lexer.Equals),
					Operands: []parser.Operand{
						lexer.Token{Text: "Name", Type: lexer.UnquotedString},
						lexer.Token{Text: "'Milk'", Type: lexer.SingleQuotedString},
					},
				},
				&parser.Operation{
					Operator: parser.Operator(lexer.LessThan),
					Operands: []parser.Operand{
						lexer.Token{Text: "Price", Type: lexer.UnquotedString},
						lexer.Token{Text: "2.55", Type: lexer.FloatingPointLiteral},
					},
				},
			},
		},
	},
	{
		input: "Name EQ 'Milk' AND Price LT 2.55",
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.And),
			Operands: []parser.Operand{
				&parser.Operation{
					Operator: parser.Operator(lexer.Equals),
					Operands: []parser.Operand{
						lexer.Token{Text: "Name", Type: lexer.UnquotedString},
						lexer.Token{Text: "'Milk'", Type: lexer.SingleQuotedString},
					},
				},
				&parser.Operation{
					Operator: parser.Operator(lexer.LessThan),
					Operands: []parser.Operand{
						lexer.Token{Text: "Price", Type: lexer.UnquotedString},
						lexer.Token{Text: "2.55", Type: lexer.FloatingPointLiteral},
					},
				},
			},
		},
	},
	{
		input: "Name eq 'Milk' AND Price LT 2.55",
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.And),
			Operands: []parser.Operand{
				&parser.Operation{
					Operator: parser.Operator(lexer.Equals),
					Operands: []parser.Operand{
						lexer.Token{Text: "Name", Type: lexer.UnquotedString},
						lexer.Token{Text: "'Milk'", Type: lexer.SingleQuotedString},
					},
				},
				&parser.Operation{
					Operator: parser.Operator(lexer.LessThan),
					Operands: []parser.Operand{
						lexer.Token{Text: "Price", Type: lexer.UnquotedString},
						lexer.Token{Text: "2.55", Type: lexer.FloatingPointLiteral},
					},
				},
			},
		},
	},
	{
		input: "Name eq 'Milk' or Price lt 2.55",
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Or),
			Operands: []parser.Operand{
				&parser.Operation{
					Operator: parser.Operator(lexer.Equals),
					Operands: []parser.Operand{
						lexer.Token{Text: "Name", Type: lexer.UnquotedString},
						lexer.Token{Text: "'Milk'", Type: lexer.SingleQuotedString},
					},
				},
				&parser.Operation{
					Operator: parser.Operator(lexer.LessThan),
					Operands: []parser.Operand{
						lexer.Token{Text: "Price", Type: lexer.UnquotedString},
						lexer.Token{Text: "2.55", Type: lexer.FloatingPointLiteral},
					},
				},
			},
		},
	},
	{
		input: "Name in ('Milk', 'Cheese')",
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.In),
			Operands: []parser.Operand{
				lexer.Token{Text: "Name", Type: lexer.UnquotedString},
				&parser.SliceOperand{Slice: []parser.Operand{
					lexer.Token{Text: "'Milk'", Type: lexer.SingleQuotedString},
					lexer.Token{Text: "'Cheese'", Type: lexer.SingleQuotedString},
				}},
			},
		},
	},
	{
		input: "Name in ['Milk', 'Cheese']",
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.In),
			Operands: []parser.Operand{
				lexer.Token{Text: "Name", Type: lexer.UnquotedString},
				&parser.SliceOperand{Slice: []parser.Operand{
					lexer.Token{Text: "'Milk'", Type: lexer.SingleQuotedString},
					lexer.Token{Text: "'Cheese'", Type: lexer.SingleQuotedString},
				}},
			},
		},
	},
	{
		input: "_id eq 6206b158000e1859781d5e16",
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Equals),
			Operands: []parser.Operand{
				lexer.Token{Text: "_id", Type: lexer.UnquotedString},
				lexer.Token{Text: "6206b158000e1859781d5e16", Type: lexer.UnquotedString},
			},
		},
	},
	{
		input: "contains(Name,'red')",
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Contains),
			Operands: []parser.Operand{
				lexer.Token{Text: "Name", Type: lexer.UnquotedString},
				lexer.Token{Text: "'red'", Type: lexer.SingleQuotedString},
			},
		},
	},
	{
		input: `Address eq {"Street":"NE 40th","City":"Redmond","State":"WA","ZipCode":"98052"}`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Equals),
			Operands: []parser.Operand{
				lexer.Token{Text: "Address", Type: lexer.UnquotedString},
				&parser.ObjectOperand{Properties: `{"Street":"NE 40th","City":"Redmond","State":"WA","ZipCode":"98052"}`},
			},
		},
	},
	{
		input: "endswith(Name,'ilk')",
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.EndsWith),
			Operands: []parser.Operand{
				lexer.Token{Text: "Name", Type: lexer.UnquotedString},
				lexer.Token{Text: "'ilk'", Type: lexer.SingleQuotedString},
			},
		},
	},
	{
		input: "not endswith(Name,'ilk')",
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Not),
			Operands: []parser.Operand{
				&parser.Operation{
					Operator: parser.Operator(lexer.EndsWith),
					Operands: []parser.Operand{
						lexer.Token{Text: "Name", Type: lexer.UnquotedString},
						lexer.Token{Text: "'ilk'", Type: lexer.SingleQuotedString},
					},
				},
			},
		},
	},
	{
		input: "length(CompanyName) eq 19",
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Equals),
			Operands: []parser.Operand{
				&parser.Operation{
					Operator: parser.Operator(lexer.Length),
					Operands: []parser.Operand{
						lexer.Token{Text: "CompanyName", Type: lexer.UnquotedString},
					},
				},
				lexer.Token{Text: "19", Type: lexer.IntegerLiteral},
			},
		},
	},
	{
		input: "startswith(CompanyName,'Futterkiste')",
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.StartsWith),
			Operands: []parser.Operand{
				lexer.Token{Text: "CompanyName", Type: lexer.UnquotedString},
				lexer.Token{Text: "'Futterkiste'", Type: lexer.SingleQuotedString},
			},
		},
	},
	{
		input: `hassubset(Names,["Milk", "Cheese"])`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.HasSubset),
			Operands: []parser.Operand{
				lexer.Token{Text: "Names", Type: lexer.UnquotedString},
				&parser.SliceOperand{Slice: []parser.Operand{
					lexer.Token{Text: `"Milk"`, Type: lexer.DoubleQuotedString},
					lexer.Token{Text: `"Cheese"`, Type: lexer.DoubleQuotedString},
				}},
			},
		},
	},
	{
		input: `Price add 2.45 eq 5.00`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Equals),
			Operands: []parser.Operand{
				&parser.Operation{
					Operator: parser.Operator(lexer.Add),
					Operands: []parser.Operand{
						lexer.Token{Text: "Price", Type: lexer.UnquotedString},
						lexer.Token{Text: "2.45", Type: lexer.FloatingPointLiteral},
					},
				},
				lexer.Token{Text: "5.00", Type: lexer.FloatingPointLiteral},
			},
		},
	},
	{
		input: `Price sub 0.55 eq 2.00`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Equals),
			Operands: []parser.Operand{
				&parser.Operation{
					Operator: parser.Operator(lexer.Subtract),
					Operands: []parser.Operand{
						lexer.Token{Text: "Price", Type: lexer.UnquotedString},
						lexer.Token{Text: "0.55", Type: lexer.FloatingPointLiteral},
					},
				},
				lexer.Token{Text: "2.00", Type: lexer.FloatingPointLiteral},
			},
		},
	},
	{
		input: `Price mul 2.0 eq 5.10`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Equals),
			Operands: []parser.Operand{
				&parser.Operation{
					Operator: parser.Operator(lexer.Multiply),
					Operands: []parser.Operand{
						lexer.Token{Text: "Price", Type: lexer.UnquotedString},
						lexer.Token{Text: "2.0", Type: lexer.FloatingPointLiteral},
					},
				},
				lexer.Token{Text: "5.10", Type: lexer.FloatingPointLiteral},
			},
		},
	},
	{
		input: `Price div 2.55 eq 1`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Equals),
			Operands: []parser.Operand{
				&parser.Operation{
					Operator: parser.Operator(lexer.Divide),
					Operands: []parser.Operand{
						lexer.Token{Text: "Price", Type: lexer.UnquotedString},
						lexer.Token{Text: "2.55", Type: lexer.FloatingPointLiteral},
					},
				},
				lexer.Token{Text: "1", Type: lexer.IntegerLiteral},
			},
		},
	},
	{
		input: `Price div 2 eq 2`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Equals),
			Operands: []parser.Operand{
				&parser.Operation{
					Operator: parser.Operator(lexer.Divide),
					Operands: []parser.Operand{
						lexer.Token{Text: "Price", Type: lexer.UnquotedString},
						lexer.Token{Text: "2", Type: lexer.IntegerLiteral},
					},
				},
				lexer.Token{Text: "2", Type: lexer.IntegerLiteral},
			},
		},
	},
	{
		input: `Price divby 2 eq 2.5`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Equals),
			Operands: []parser.Operand{
				&parser.Operation{
					Operator: parser.Operator(lexer.DivideFloat),
					Operands: []parser.Operand{
						lexer.Token{Text: "Price", Type: lexer.UnquotedString},
						lexer.Token{Text: "2", Type: lexer.IntegerLiteral},
					},
				},
				lexer.Token{Text: "2.5", Type: lexer.FloatingPointLiteral},
			},
		},
	},
	{
		input: `Rating mod 5 eq 0`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Equals),
			Operands: []parser.Operand{
				&parser.Operation{
					Operator: parser.Operator(lexer.Modulo),
					Operands: []parser.Operand{
						lexer.Token{Text: "Rating", Type: lexer.UnquotedString},
						lexer.Token{Text: "5", Type: lexer.IntegerLiteral},
					},
				},
				lexer.Token{Text: "0", Type: lexer.IntegerLiteral},
			},
		},
	},
	{
		input: `(4 add 5) mod (4 sub 1) eq 0`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Equals),
			Operands: []parser.Operand{
				&parser.Operation{
					Operator: parser.Operator(lexer.Modulo),
					Operands: []parser.Operand{
						&parser.Operation{
							Operator: parser.Operator(lexer.Add),
							Operands: []parser.Operand{
								lexer.Token{Text: "4", Type: lexer.IntegerLiteral},
								lexer.Token{Text: "5", Type: lexer.IntegerLiteral},
							},
						},
						&parser.Operation{
							Operator: parser.Operator(lexer.Subtract),
							Operands: []parser.Operand{
								lexer.Token{Text: "4", Type: lexer.IntegerLiteral},
								lexer.Token{Text: "1", Type: lexer.IntegerLiteral},
							},
						},
					},
				},
				lexer.Token{Text: "0", Type: lexer.IntegerLiteral},
			},
		},
	},
	{
		input: `concat(concat(City,', '),Country) eq 'Berlin, Germany'`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Equals),
			Operands: []parser.Operand{
				&parser.Operation{
					Operator: parser.Operator(lexer.Concat),
					Operands: []parser.Operand{
						&parser.Operation{
							Operator: parser.Operator(lexer.Concat),
							Operands: []parser.Operand{
								lexer.Token{Text: "City", Type: lexer.UnquotedString},
								lexer.Token{Text: "', '", Type: lexer.SingleQuotedString},
							},
						},
						lexer.Token{Text: "Country", Type: lexer.UnquotedString},
					},
				},
				lexer.Token{Text: "'Berlin, Germany'", Type: lexer.SingleQuotedString},
			},
		},
	},
	{
		input: `indexof(CompanyName,'lfreds') eq 1`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Equals),
			Operands: []parser.Operand{
				&parser.Operation{
					Operator: parser.Operator(lexer.IndexOf),
					Operands: []parser.Operand{
						lexer.Token{Text: "CompanyName", Type: lexer.UnquotedString},
						lexer.Token{Text: "'lfreds'", Type: lexer.SingleQuotedString},
					},
				},
				lexer.Token{Text: "1", Type: lexer.IntegerLiteral},
			},
		},
	},
	{
		input: `substring(CompanyName,1) eq 'lfreds Futterkiste'`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Equals),
			Operands: []parser.Operand{
				&parser.Operation{
					Operator: parser.Operator(lexer.Substring),
					Operands: []parser.Operand{
						lexer.Token{Text: "CompanyName", Type: lexer.UnquotedString},
						lexer.Token{Text: "1", Type: lexer.IntegerLiteral},
					},
				},
				lexer.Token{Text: "'lfreds Futterkiste'", Type: lexer.SingleQuotedString},
			},
		},
	},
	{
		input: `substring(CompanyName,1,2) eq 'lf'`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Equals),
			Operands: []parser.Operand{
				&parser.Operation{
					Operator: parser.Operator(lexer.Substring),
					Operands: []parser.Operand{
						lexer.Token{Text: "CompanyName", Type: lexer.UnquotedString},
						lexer.Token{Text: "1", Type: lexer.IntegerLiteral},
						lexer.Token{Text: "2", Type: lexer.IntegerLiteral},
					},
				},
				lexer.Token{Text: "'lf'", Type: lexer.SingleQuotedString},
			},
		},
	},
	{
		input: `hassubsequence([4,1,3],[4,1,3])`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.HasSubsequence),
			Operands: []parser.Operand{
				&parser.SliceOperand{Slice: []parser.Operand{
					lexer.Token{Text: "4", Type: lexer.IntegerLiteral},
					lexer.Token{Text: "1", Type: lexer.IntegerLiteral},
					lexer.Token{Text: "3", Type: lexer.IntegerLiteral},
				}},
				&parser.SliceOperand{Slice: []parser.Operand{
					lexer.Token{Text: "4", Type: lexer.IntegerLiteral},
					lexer.Token{Text: "1", Type: lexer.IntegerLiteral},
					lexer.Token{Text: "3", Type: lexer.IntegerLiteral},
				}},
			},
		},
	},
	{
		input: `hassubsequence([4,1,3],[4,1])`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.HasSubsequence),
			Operands: []parser.Operand{
				&parser.SliceOperand{Slice: []parser.Operand{
					lexer.Token{Text: "4", Type: lexer.IntegerLiteral},
					lexer.Token{Text: "1", Type: lexer.IntegerLiteral},
					lexer.Token{Text: "3", Type: lexer.IntegerLiteral},
				}},
				&parser.SliceOperand{Slice: []parser.Operand{
					lexer.Token{Text: "4", Type: lexer.IntegerLiteral},
					lexer.Token{Text: "1", Type: lexer.IntegerLiteral},
				}},
			},
		},
	},
	{
		input: `hassubsequence([4,1,3],[4,3])`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.HasSubsequence),
			Operands: []parser.Operand{
				&parser.SliceOperand{Slice: []parser.Operand{
					lexer.Token{Text: "4", Type: lexer.IntegerLiteral},
					lexer.Token{Text: "1", Type: lexer.IntegerLiteral},
					lexer.Token{Text: "3", Type: lexer.IntegerLiteral},
				}},
				&parser.SliceOperand{Slice: []parser.Operand{
					lexer.Token{Text: "4", Type: lexer.IntegerLiteral},
					lexer.Token{Text: "3", Type: lexer.IntegerLiteral},
				}},
			},
		},
	},
	{
		input: `hassubsequence([4,1,3,1],[1,1])`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.HasSubsequence),
			Operands: []parser.Operand{
				&parser.SliceOperand{Slice: []parser.Operand{
					lexer.Token{Text: "4", Type: lexer.IntegerLiteral},
					lexer.Token{Text: "1", Type: lexer.IntegerLiteral},
					lexer.Token{Text: "3", Type: lexer.IntegerLiteral},
					lexer.Token{Text: "1", Type: lexer.IntegerLiteral},
				}},
				&parser.SliceOperand{Slice: []parser.Operand{
					lexer.Token{Text: "1", Type: lexer.IntegerLiteral},
					lexer.Token{Text: "1", Type: lexer.IntegerLiteral},
				}},
			},
		},
	},
	{
		input: `hassubsequence([4,1,3],[1,3,4])`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.HasSubsequence),
			Operands: []parser.Operand{
				&parser.SliceOperand{Slice: []parser.Operand{
					lexer.Token{Text: "4", Type: lexer.IntegerLiteral},
					lexer.Token{Text: "1", Type: lexer.IntegerLiteral},
					lexer.Token{Text: "3", Type: lexer.IntegerLiteral},
				}},
				&parser.SliceOperand{Slice: []parser.Operand{
					lexer.Token{Text: "1", Type: lexer.IntegerLiteral},
					lexer.Token{Text: "3", Type: lexer.IntegerLiteral},
					lexer.Token{Text: "4", Type: lexer.IntegerLiteral},
				}},
			},
		},
	},
	{
		input: `hassubsequence([4,1,3],[3,1])`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.HasSubsequence),
			Operands: []parser.Operand{
				&parser.SliceOperand{Slice: []parser.Operand{
					lexer.Token{Text: "4", Type: lexer.IntegerLiteral},
					lexer.Token{Text: "1", Type: lexer.IntegerLiteral},
					lexer.Token{Text: "3", Type: lexer.IntegerLiteral},
				}},
				&parser.SliceOperand{Slice: []parser.Operand{
					lexer.Token{Text: "3", Type: lexer.IntegerLiteral},
					lexer.Token{Text: "1", Type: lexer.IntegerLiteral},
				}},
			},
		},
	},
	{
		input: `hassubsequence([1,2],[1,1,2])`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.HasSubsequence),
			Operands: []parser.Operand{
				&parser.SliceOperand{Slice: []parser.Operand{
					lexer.Token{Text: "1", Type: lexer.IntegerLiteral},
					lexer.Token{Text: "2", Type: lexer.IntegerLiteral},
				}},
				&parser.SliceOperand{Slice: []parser.Operand{
					lexer.Token{Text: "1", Type: lexer.IntegerLiteral},
					lexer.Token{Text: "1", Type: lexer.IntegerLiteral},
					lexer.Token{Text: "2", Type: lexer.IntegerLiteral},
				}},
			},
		},
	},
	{
		input: `matchesPattern(CompanyName,'%5EA.*e$')`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.MatchesPattern),
			Operands: []parser.Operand{
				lexer.Token{Text: "CompanyName", Type: lexer.UnquotedString},
				lexer.Token{Text: "'%5EA.*e$'", Type: lexer.SingleQuotedString},
			},
		},
	},
	{
		input: `tolower(CompanyName) eq 'alfreds futterkiste'`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Equals),
			Operands: []parser.Operand{
				&parser.Operation{
					Operator: parser.Operator(lexer.ToLower),
					Operands: []parser.Operand{
						lexer.Token{Text: "CompanyName", Type: lexer.UnquotedString},
					},
				},
				lexer.Token{Text: "'alfreds futterkiste'", Type: lexer.SingleQuotedString},
			},
		},
	},
	{
		input: `toupper(CompanyName) eq 'ALFREDS FUTTERKISTE'`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Equals),
			Operands: []parser.Operand{
				&parser.Operation{
					Operator: parser.Operator(lexer.ToUpper),
					Operands: []parser.Operand{
						lexer.Token{Text: "CompanyName", Type: lexer.UnquotedString},
					},
				},
				lexer.Token{Text: "'ALFREDS FUTTERKISTE'", Type: lexer.SingleQuotedString},
			},
		},
	},
	{
		input: `trim(CompanyName) eq CompanyName`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Equals),
			Operands: []parser.Operand{
				&parser.Operation{
					Operator: parser.Operator(lexer.Trim),
					Operands: []parser.Operand{
						lexer.Token{Text: "CompanyName", Type: lexer.UnquotedString},
					},
				},
				lexer.Token{Text: "CompanyName", Type: lexer.UnquotedString},
			},
		},
	},
	{
		input: `day(BirthDate) eq 8`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Equals),
			Operands: []parser.Operand{
				&parser.Operation{
					Operator: parser.Operator(lexer.Day),
					Operands: []parser.Operand{
						lexer.Token{Text: "BirthDate", Type: lexer.UnquotedString},
					},
				},
				lexer.Token{Text: "8", Type: lexer.IntegerLiteral},
			},
		},
	},
	{
		input: `fractionalseconds(BirthDate) lt 0.1`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.LessThan),
			Operands: []parser.Operand{
				&parser.Operation{
					Operator: parser.Operator(lexer.FractionalSeconds),
					Operands: []parser.Operand{
						lexer.Token{Text: "BirthDate", Type: lexer.UnquotedString},
					},
				},
				lexer.Token{Text: "0.1", Type: lexer.FloatingPointLiteral},
			},
		},
	},
	{
		input: `hour(BirthDate) eq 4`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Equals),
			Operands: []parser.Operand{
				&parser.Operation{
					Operator: parser.Operator(lexer.Hour),
					Operands: []parser.Operand{
						lexer.Token{Text: "BirthDate", Type: lexer.UnquotedString},
					},
				},
				lexer.Token{Text: "4", Type: lexer.IntegerLiteral},
			},
		},
	},
	{
		input: `minute(BirthDate) eq 40`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Equals),
			Operands: []parser.Operand{
				&parser.Operation{
					Operator: parser.Operator(lexer.Minute),
					Operands: []parser.Operand{
						lexer.Token{Text: "BirthDate", Type: lexer.UnquotedString},
					},
				},
				lexer.Token{Text: "40", Type: lexer.IntegerLiteral},
			},
		},
	},
	{
		input: `month(BirthDate) eq 5`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Equals),
			Operands: []parser.Operand{
				&parser.Operation{
					Operator: parser.Operator(lexer.Month),
					Operands: []parser.Operand{
						lexer.Token{Text: "BirthDate", Type: lexer.UnquotedString},
					},
				},
				lexer.Token{Text: "5", Type: lexer.IntegerLiteral},
			},
		},
	},
	{
		input: `second(BirthDate) eq 40`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Equals),
			Operands: []parser.Operand{
				&parser.Operation{
					Operator: parser.Operator(lexer.Second),
					Operands: []parser.Operand{
						lexer.Token{Text: "BirthDate", Type: lexer.UnquotedString},
					},
				},
				lexer.Token{Text: "40", Type: lexer.IntegerLiteral},
			},
		},
	},
	{
		input: `year(BirthDate) eq 1971`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Equals),
			Operands: []parser.Operand{
				&parser.Operation{
					Operator: parser.Operator(lexer.Year),
					Operands: []parser.Operand{
						lexer.Token{Text: "BirthDate", Type: lexer.UnquotedString},
					},
				},
				lexer.Token{Text: "1971", Type: lexer.IntegerLiteral},
			},
		},
	},
	{
		input: `ceiling(Freight) eq 32`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Equals),
			Operands: []parser.Operand{
				&parser.Operation{
					Operator: parser.Operator(lexer.Ceiling),
					Operands: []parser.Operand{
						lexer.Token{Text: "Freight", Type: lexer.UnquotedString},
					},
				},
				lexer.Token{Text: "32", Type: lexer.IntegerLiteral},
			},
		},
	},
	{
		input: `floor(Freight) eq 32`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Equals),
			Operands: []parser.Operand{
				&parser.Operation{
					Operator: parser.Operator(lexer.Floor),
					Operands: []parser.Operand{
						lexer.Token{Text: "Freight", Type: lexer.UnquotedString},
					},
				},
				lexer.Token{Text: "32", Type: lexer.IntegerLiteral},
			},
		},
	},
	{
		input: `round(Freight) eq 32`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Equals),
			Operands: []parser.Operand{
				&parser.Operation{
					Operator: parser.Operator(lexer.Round),
					Operands: []parser.Operand{
						lexer.Token{Text: "Freight", Type: lexer.UnquotedString},
					},
				},
				lexer.Token{Text: "32", Type: lexer.IntegerLiteral},
			},
		},
	},
	{
		input: `DiscontinuedDate eq null`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Equals),
			Operands: []parser.Operand{
				lexer.Token{Text: "DiscontinuedDate", Type: lexer.UnquotedString},
				lexer.Token{Text: "null", Type: lexer.NullLiteral},
			},
		},
	},
	{
		input: `year eq ':0'`,
		expectedOperation: parser.Operation{
			Operator: parser.Operator(lexer.Equals),
			Operands: []parser.Operand{
				lexer.Token{Text: "year", Type: lexer.UnquotedString},
				lexer.Token{Text: "':0'", Type: lexer.SingleQuotedString},
			},
		},
	},
}

func TestGetExpression(t *testing.T) {
	t.Parallel()
	for _, test := range testCases {
		tc := test
		t.Run(test.input, func(t *testing.T) {
			t.Parallel()
			parser, err := parser.NewParser(tc.input)
			if err != nil {
				t.Errorf("error creating parser: %s", err)
			}
			op, err := parser.GetOperation()
			if err != nil {
				t.Fatalf("error getting operation: %s", err)
			}
			err = compareOperations(op, &tc.expectedOperation)
			if err != nil {
				json1Got, _ := json.Marshal(op)
				json2Got, _ := json.Marshal(tc.expectedOperation)
				t.Errorf("Expression %s\nparsed to: %s\nexpected: %s\nError: %v", tc.input, string(json1Got), string(json2Got), err)
			}
		})
	}
}

func BenchmarkGetExpression(b *testing.B) {
	for _, test := range testCases {
		tc := test
		b.Run(test.input, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				parser, err := parser.NewParser(tc.input)
				if err != nil {
					b.Errorf("error creating parser: %s", err)
				}
				_, err = parser.GetOperation()
				if err != nil {
					b.Fatalf("error getting operation: %s", err)
				}
			}
		})
	}
}

type TestError struct {
	message string
}

func (e *TestError) Error() string {
	return e.message
}

func newTestError(format string, a ...interface{}) error {
	return &TestError{message: fmt.Sprintf(format, a...)}
}

func compareOperations(got, expected *parser.Operation) error {
	if expected == nil {
		return nil
	}
	if got == nil {
		return newTestError("got nil operation")
	}
	if int(got.Operator) != int(expected.Operator) {
		return newTestError("got different operators %d != %d", got.Operator, expected.Operator)
	}
	if len(got.Operands) != len(expected.Operands) {
		return newTestError("got different operand lengths %d != %d", len(got.Operands), len(expected.Operands))
	}
	for i, operand := range got.Operands {
		err := compareOperand(operand, expected.Operands[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func compareOperand(got, expected parser.Operand) error {
	gotData, err := got.GetData()
	if err != nil {
		return newTestError("error getting data: %s", err)
	}
	expectedData, err := expected.GetData()
	if err != nil {
		return newTestError("error getting expected data: %s", err)
	}
	if reflect.TypeOf(gotData) != reflect.TypeOf(expectedData) {
		return newTestError("got different types %T != %T", gotData, expectedData)
	}
	return compareData(gotData, expectedData)
}

//nolint:cyclop,funlen,gocognit
func compareData(gotData, expectedData interface{}) error {
	switch gotType := gotData.(type) {
	case *parser.Operation:
		expectedOp, ok := expectedData.(*parser.Operation)
		if !ok {
			return newTestError("expected data is not an operation")
		}
		err := compareOperations(gotType, expectedOp)
		if err != nil {
			return err
		}
	case []*lexer.Token:
		expectedSlice, ok := expectedData.([]*lexer.Token)
		if !ok {
			return newTestError("expected data is not a slice of tokens %T", expectedData)
		}
		if len(gotType) != len(expectedSlice) {
			return newTestError("got different slice lengths %d != %d", len(gotType), len(expectedSlice))
		}
		for j, token := range gotType {
			err := tokenCompare(token, expectedSlice[j])
			if err != nil {
				return err
			}
		}
	case string:
		expectedString, ok := expectedData.(string)
		if !ok {
			return newTestError("expected data is not a string while got a string")
		}
		if gotType != expectedString {
			return newTestError("%s != %s", gotType, expectedString)
		}
	case []parser.Operand:
		for j, operand := range gotType {
			expectedArr, ok := expectedData.([]parser.Operand)
			if !ok {
				return newTestError("expected data is not a slice of operands %T", expectedData)
			}
			err := compareOperandData(operand, expectedArr[j])
			if err != nil {
				return err
			}
		}
	case map[string]interface{}:
		expectedMap, ok := expectedData.(map[string]interface{})
		if !ok {
			return newTestError("expected data is not a map %T", expectedData)
		}
		err := mapCompare(gotType, expectedMap)
		if err != nil {
			return err
		}
	default:
		if gotData != expectedData {
			return newTestError("%#v != %#v", gotData, expectedData)
		}
	}
	return nil
}

func compareOperandData(got, expected parser.Operand) error {
	switch myOperand := got.(type) {
	case *parser.Operation:
		expectedOp, ok := expected.(*parser.Operation)
		if !ok {
			return newTestError("expected data is not an operation %T", expected)
		}
		err := compareOperations(myOperand, expectedOp)
		if err != nil {
			return err
		}
	case *lexer.Token:
		expectedToken, ok := expected.(*lexer.Token)
		if !ok {
			myTok, ok := expected.(lexer.Token)
			if !ok {
				return newTestError("expected data is not a token or *token %T", expected)
			}
			expectedToken = &myTok
		}
		err := tokenCompare(myOperand, expectedToken)
		if err != nil {
			return err
		}
	default:
		if got != expected {
			return newTestError("operand compare failed %T\n", got)
		}
	}
	return nil
}

func tokenCompare(got, expected *lexer.Token) error {
	if got.Type != expected.Type {
		return newTestError("got different types %d != %d", got.Type, expected.Type)
	}
	gotData, err := got.GetData()
	if err != nil {
		return newTestError("error getting data: %v", err)
	}
	expectedData, err := expected.GetData()
	if err != nil {
		return newTestError("error getting expected data: %v", err)
	}
	if gotData != expectedData {
		return newTestError("%#v != %#v", gotData, expectedData)
	}
	return nil
}

func mapCompare(got, expected map[string]interface{}) error {
	//nolint:errchkjson
	json1Got, _ := json.Marshal(got)
	//nolint:errchkjson
	json2Got, _ := json.Marshal(expected)
	if string(json1Got) != string(json2Got) {
		return newTestError("got different maps %s != %s", string(json1Got), string(json2Got))
	}
	return nil
}
