package main

import (
	"os"

	"github.com/progrium/macdriver/gen"
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "FIXME",
	// Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pkg := gen.GoPackage{
			Package: "foundation",

			Imports: []gen.Import{
				{Path: "unsafe"},
				{Path: "github.com/progrium/macdriver/objc2", Alias: "objc"},
			},

			MsgSendWrappers: []gen.CGoMsgSend{
				{
					Name:        "NSString_init",
					Selector:    "init",
					MsgSendFunc: "objc_msgSend",
					Args:        []gen.Arg{},
					Return:      "void*",
				},
				{
					Name:        "NSString_initWithBytes_length_encoding_",
					Selector:    "initWithBytes:length:encoding:",
					MsgSendFunc: "objc_msgSend",
					Args: []gen.Arg{
						{Name: "bytes", Type: "void*"},
						{Name: "len", Type: "unsigned long"},
						{Name: "encoding", Type: "unsigned long"},
					},
					Return: "void*",
				},
				{
					Name:        "NSString_length",
					Selector:    "length",
					MsgSendFunc: "objc_msgSend",
					Args:        []gen.Arg{},
					Return:      "unsigned long",
				},
				{
					Name:        "NSString_characterAtIndex_",
					Selector:    "characterAtIndex:",
					MsgSendFunc: "objc_msgSend",
					Args: []gen.Arg{
						{Name: "index", Type: "unsigned long"},
					},
					Return: "unsigned short",
				},
			},

			CGoWrapperFuncs: []gen.CGoWrapperFunc{
				{
					Name: "NSString_init",
					Args: []gen.CGoWrapperArg{
						{Name: "self", Type: "objc.Object", ToCGoFmt: "unsafe.Pointer(%s.Id())"},
					},
					Returns: []gen.CGoWrapperReturn{
						{Type: "NSString", FromCGoFmt: "NSString_fromPointer(%s)"},
					},
				},
				{
					Name: "NSString_initWithBytes_length_encoding_",
					Args: []gen.CGoWrapperArg{
						{Name: "self", Type: "objc.Object", ToCGoFmt: "unsafe.Pointer(%s.Id())"},
						{Name: "bytes", Type: "unsafe.Pointer", ToCGoFmt: "%s"},
						{Name: "len", Type: "objc.NSUInteger", ToCGoFmt: "C.ulong(%s)"},
						{Name: "encoding", Type: "objc.NSStringEncoding", ToCGoFmt: "C.ulong(%s)"},
					},
					Returns: []gen.CGoWrapperReturn{
						{Type: "NSString", FromCGoFmt: "NSString_fromPointer(%s)"},
					},
				},
				{
					Name: "NSString_length",
					Args: []gen.CGoWrapperArg{
						{Name: "self", Type: "objc.Object", ToCGoFmt: "unsafe.Pointer(%s.Id())"},
					},
					Returns: []gen.CGoWrapperReturn{
						{Type: "objc.NSUInteger", FromCGoFmt: "objc.NSUInteger(%s)"},
					},
				},
				{
					Name: "NSString_characterAtIndex_",
					Args: []gen.CGoWrapperArg{
						{Name: "self", Type: "objc.Object", ToCGoFmt: "unsafe.Pointer(%s.Id())"},
						{Name: "index", Type: "objc.NSUInteger", ToCGoFmt: "C.ulong(%s)"},
					},
					Returns: []gen.CGoWrapperReturn{
						{Type: "objc.Unichar", FromCGoFmt: "objc.Unichar(%s)"},
					},
				},
			},

			Classes: []gen.ClassDef{
				{
					Name: "NSString",
					InstanceMethods: []gen.MethodDef{
						{
							Name:        "Init",
							WrappedFunc: "NSString_init",
							Args:        []gen.Arg{},
							Returns: []gen.MethodReturn{
								{Type: "NSString"},
							},
						},
						{
							Name:        "InitWithBytes_length_encoding_",
							WrappedFunc: "NSString_initWithBytes_length_encoding_",
							Args: []gen.Arg{
								{Name: "bytes", Type: "unsafe.Pointer"},
								{Name: "len", Type: "objc.NSUInteger"},
								{Name: "encoding", Type: "objc.NSStringEncoding"},
							},
							Returns: []gen.MethodReturn{
								{Type: "NSString"},
							},
						},
						{
							Name:        "Length",
							WrappedFunc: "NSString_length",
							Args:        []gen.Arg{},
							Returns: []gen.MethodReturn{
								{Type: "objc.NSUInteger"},
							},
						},
						{
							Name:        "CharacterAtIndex_",
							WrappedFunc: "NSString_characterAtIndex_",
							Args: []gen.Arg{
								{Name: "index", Type: "objc.NSUInteger"},
							},
							Returns: []gen.MethodReturn{
								{Type: "objc.Unichar"},
							},
						},
					},
				},
			},
		}
		fatal(pkg.Generate(os.Stdout))
	},
}
