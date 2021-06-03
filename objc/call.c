// #cgo LDFLAGS: -lobjc
#define __OBJC2__ 1
#include <objc/message.h>
#include "_cgo_export.h"

void GoObjc_Invoke(id self, SEL _cmd, id invocation) {
	goInvoke((void*)self, (void*)_cmd, (void*)invocation);
}

id GoObjc_MethodSignature(id self, SEL _cmd, SEL sel) {
	return goMethodSignature((void*)self, (void*)_cmd, (void*)sel);
}
