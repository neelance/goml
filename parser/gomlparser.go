package parser

import (
	"go/ast"
	"go/token"
)

type Mode uint

const (
	PackageClauseOnly Mode             = 1 << iota // stop parsing after package clause
	ImportsOnly                                    // stop parsing after import declarations
	ParseComments                                  // parse comments and add them to AST
	Trace                                          // print a trace of parsed productions
	DeclarationErrors                              // report declaration errors
	SpuriousErrors                                 // same as AllErrors, for backward-compatibility
	AllErrors         = SpuriousErrors             // report all errors (not just the first 10 on different lines)
)

var nilIdent = ast.NewIdent("nil")
var trueIdent = ast.NewIdent("true")
var elementType = ast.NewIdent("element")
var attrsType = ast.NewIdent("attributes")

func ParseFile(fset *token.FileSet, filename string, src []byte) (f *ast.File, err error) {
	var p parser
	defer func() {
		if e := recover(); e != nil {
			// resume same panic if it's not a bailout
			if _, ok := e.(bailout); !ok {
				panic(e)
			}
		}

		p.errors.Sort()
		err = p.errors.Err()
	}()

	p.init(fset, filename, src, ParseComments)
	f = p.parseFile()
	return
}

func (p *parser) parseTag() ast.Stmt {
	if p.trace {
		defer un(trace(p, "Tag"))
	}

	pos := p.expect(token.LSS)
	name := p.parseIdent()

	var attrs ast.Expr = nilIdent
	if p.tok == token.LPAREN {
		p.expect(token.LPAREN)
		var elts []ast.Expr
		for {
			key := ""
			if p.tok == token.PERIOD {
				p.expect(token.PERIOD)
				key += "."
			}
			key += p.parseIdent().Name

			var value ast.Expr = trueIdent
			if p.tok == token.COLON {
				p.expect(token.COLON)
				value = p.parseRhs()
			}

			elts = append(elts, &ast.KeyValueExpr{
				Key: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"` + key + `"`,
				},
				Value: value,
			})

			if p.tok != token.COMMA {
				break
			}
			p.expect(token.COMMA)
		}
		p.expect(token.RPAREN)
		attrs = &ast.CompositeLit{
			Type: attrsType,
			Elts: elts,
		}
	}

	p.expect(token.GTR)

	var body ast.Expr = nilIdent
	if p.tok == token.LBRACE {
		block := p.parseBlockStmt()
		body = &ast.FuncLit{
			Type: &ast.FuncType{
				Func: block.Lbrace,
				Params: &ast.FieldList{
					List: []*ast.Field{
						&ast.Field{
							Names: []*ast.Ident{ast.NewIdent("e")},
							Type:  elementType,
						},
					},
				},
			},
			Body: block,
		}
	}

	return &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   &ast.Ident{NamePos: pos, Name: "e"},
				Sel: ast.NewIdent("AppendElement"),
			},
			Args: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: `"` + name.Name + `"`,
				},
				attrs,
				body,
			},
		},
	}
}

func (p *parser) parseTextNode() ast.Stmt {
	pos := p.expect(token.REM)
	content := p.parseRhs()

	return &ast.ExprStmt{
		X: &ast.CallExpr{
			Lparen: content.Pos(),
			Rparen: content.End(),
			Fun: &ast.SelectorExpr{
				X:   &ast.Ident{NamePos: pos, Name: "e"},
				Sel: ast.NewIdent("AppendTextNode"),
			},
			Args: []ast.Expr{
				content,
			},
		},
	}
}
