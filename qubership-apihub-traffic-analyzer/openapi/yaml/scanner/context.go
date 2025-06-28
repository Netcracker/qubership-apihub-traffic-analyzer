// Copyright 2024-2025 NetCracker Technology Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package scanner

import (
	"sync"

	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/openapi/yaml/token"
)

// Context context at scanning
type Context struct {
	idx                int
	size               int
	notSpaceCharPos    int
	notSpaceOrgCharPos int
	src                []rune
	buf                []rune
	oBuf               []rune
	tokens             token.Tokens
	isRawFolded        bool
	isLiteral          bool
	isFolded           bool
	isSingleLine       bool
	literalOpt         string
}

var (
	ctxPool = sync.Pool{
		New: func() interface{} {
			return createContext()
		},
	}
)

func createContext() *Context {
	return &Context{
		idx:          0,
		tokens:       token.Tokens{},
		isSingleLine: true,
	}
}

func newContext(src []rune) *Context {
	ctx := ctxPool.Get().(*Context)
	ctx.reset(src)
	return ctx
}

func (c *Context) release() {
	ctxPool.Put(c)
}

func (c *Context) reset(src []rune) {
	c.idx = 0
	c.size = len(src)
	c.src = src
	c.tokens = c.tokens[:0]
	c.resetBuffer()
	c.isSingleLine = true
}

func (c *Context) resetBuffer() {
	c.buf = c.buf[:0]
	c.oBuf = c.oBuf[:0]
	c.notSpaceCharPos = 0
	c.notSpaceOrgCharPos = 0
}

func (c *Context) isSaveIndentMode() bool {
	return c.isLiteral || c.isFolded || c.isRawFolded
}

func (c *Context) breakLiteral() {
	c.isLiteral = false
	c.isRawFolded = false
	c.isFolded = false
	c.literalOpt = ""
}

func (c *Context) addToken(tk *token.Token) {
	if tk == nil {
		return
	}
	c.tokens = append(c.tokens, tk)
}

func (c *Context) addBuf(r rune) {
	if len(c.buf) == 0 && r == ' ' {
		return
	}
	c.buf = append(c.buf, r)
	if r != ' ' && r != '\t' {
		c.notSpaceCharPos = len(c.buf)
	}
}

func (c *Context) addOriginBuf(r rune) {
	c.oBuf = append(c.oBuf, r)
	if r != ' ' && r != '\t' {
		c.notSpaceOrgCharPos = len(c.oBuf)
	}
}

func (c *Context) removeRightSpaceFromBuf() int {
	trimmedBuf := c.oBuf[:c.notSpaceOrgCharPos]
	bufLen := len(trimmedBuf)
	diff := len(c.oBuf) - bufLen
	if diff > 0 {
		c.oBuf = c.oBuf[:bufLen]
		c.buf = c.bufferedSrc()
	}
	return diff
}

func (c *Context) isDocument() bool {
	return c.isLiteral || c.isFolded || c.isRawFolded
}

func (c *Context) isEOS() bool {
	return len(c.src)-1 <= c.idx
}

func (c *Context) isNextEOS() bool {
	return len(c.src)-1 <= c.idx+1
}

func (c *Context) next() bool {
	return c.idx < c.size
}

func (c *Context) source(s, e int) string {
	return string(c.src[s:e])
}

func (c *Context) previousChar() rune {
	if c.idx > 0 {
		return c.src[c.idx-1]
	}
	return rune(0)
}

func (c *Context) currentChar() rune {
	return c.src[c.idx]
}

func (c *Context) nextChar() rune {
	if c.size > c.idx+1 {
		return c.src[c.idx+1]
	}
	return rune(0)
}

func (c *Context) repeatNum(r rune) int {
	cnt := 0
	for i := c.idx; i < c.size; i++ {
		if c.src[i] == r {
			cnt++
		} else {
			break
		}
	}
	return cnt
}

func (c *Context) progress(num int) {
	c.idx += num
}

func (c *Context) nextPos() int {
	return c.idx + 1
}

func (c *Context) existsBuffer() bool {
	return len(c.bufferedSrc()) != 0
}

func (c *Context) bufferedSrc() []rune {
	src := c.buf[:c.notSpaceCharPos]
	if len(src) > 0 && src[len(src)-1] == '\n' && c.isDocument() && c.literalOpt == "-" {
		// remove end '\n' character
		src = src[:len(src)-1]
	}
	return src
}

func (c *Context) bufferedToken(pos *token.Position) *token.Token {
	if c.idx == 0 {
		return nil
	}
	source := c.bufferedSrc()
	if len(source) == 0 {
		return nil
	}
	var tk *token.Token
	if c.isDocument() {
		tk = token.String(string(source), string(c.oBuf), pos)
	} else {
		tk = token.New(string(source), string(c.oBuf), pos)
	}
	c.resetBuffer()
	return tk
}
