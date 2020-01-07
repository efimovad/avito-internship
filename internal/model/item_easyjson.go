// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package model

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonA80d3b19DecodeGithubComEfimovadAvitoInternshipInternalModel(in *jlexer.Lexer, out *Params) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "Date":
			out.Date = bool(in.Bool())
		case "Price":
			out.Price = bool(in.Bool())
		case "Desc":
			out.Desc = bool(in.Bool())
		case "Page":
			out.Page = int64(in.Int64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonA80d3b19EncodeGithubComEfimovadAvitoInternshipInternalModel(out *jwriter.Writer, in Params) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Date\":"
		out.RawString(prefix[1:])
		out.Bool(bool(in.Date))
	}
	{
		const prefix string = ",\"Price\":"
		out.RawString(prefix)
		out.Bool(bool(in.Price))
	}
	{
		const prefix string = ",\"Desc\":"
		out.RawString(prefix)
		out.Bool(bool(in.Desc))
	}
	{
		const prefix string = ",\"Page\":"
		out.RawString(prefix)
		out.Int64(int64(in.Page))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Params) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonA80d3b19EncodeGithubComEfimovadAvitoInternshipInternalModel(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Params) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonA80d3b19EncodeGithubComEfimovadAvitoInternshipInternalModel(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Params) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonA80d3b19DecodeGithubComEfimovadAvitoInternshipInternalModel(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Params) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonA80d3b19DecodeGithubComEfimovadAvitoInternshipInternalModel(l, v)
}
func easyjsonA80d3b19DecodeGithubComEfimovadAvitoInternshipInternalModel1(in *jlexer.Lexer, out *Item) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "title":
			out.Title = string(in.String())
		case "description":
			out.Description = string(in.String())
		case "price":
			out.Price = float32(in.Float32())
		case "mainImage":
			out.MainImage = string(in.String())
		case "images":
			if in.IsNull() {
				in.Skip()
				out.Images = nil
			} else {
				in.Delim('[')
				if out.Images == nil {
					if !in.IsDelim(']') {
						out.Images = make([]string, 0, 4)
					} else {
						out.Images = []string{}
					}
				} else {
					out.Images = (out.Images)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Images = append(out.Images, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonA80d3b19EncodeGithubComEfimovadAvitoInternshipInternalModel1(out *jwriter.Writer, in Item) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"title\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Title))
	}
	if in.Description != "" {
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"price\":"
		out.RawString(prefix)
		out.Float32(float32(in.Price))
	}
	if in.MainImage != "" {
		const prefix string = ",\"mainImage\":"
		out.RawString(prefix)
		out.String(string(in.MainImage))
	}
	if len(in.Images) != 0 {
		const prefix string = ",\"images\":"
		out.RawString(prefix)
		{
			out.RawByte('[')
			for v2, v3 := range in.Images {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.String(string(v3))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Item) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonA80d3b19EncodeGithubComEfimovadAvitoInternshipInternalModel1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Item) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonA80d3b19EncodeGithubComEfimovadAvitoInternshipInternalModel1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Item) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonA80d3b19DecodeGithubComEfimovadAvitoInternshipInternalModel1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Item) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonA80d3b19DecodeGithubComEfimovadAvitoInternshipInternalModel1(l, v)
}
