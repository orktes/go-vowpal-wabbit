#include <stdlib.h>
#include <vowpalwabbit/vwdll.h>

#ifdef __cplusplus
extern "C" {
#endif
    VW_DLL_MEMBER void VW_CALLING_CONV VW_ReadDSJSONExample(VW_HANDLE handle, const char* line, VW_EXAMPLE* examples, size_t* example_count);
#ifdef __cplusplus
}
#endif