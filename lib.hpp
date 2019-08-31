
#include <vowpalwabbit/vwdll.h>

#ifdef __cplusplus
extern "C"
{
#endif

    VW_DLL_MEMBER VW_EXAMPLE VW_CALLING_CONV VW_ReadDSJSONExample(VW_HANDLE handle, const char* line, size_t* example_count);
    VW_DLL_MEMBER VW_EXAMPLE VW_CALLING_CONV VW_ReadJSONExample(VW_HANDLE handle, const char* line, size_t* example_count);
    VW_DLL_MEMBER void VW_CALLING_CONV VW_MultiLineLearn(VW_HANDLE handle, VW_EXAMPLE* example_handles, size_t example_count);
    VW_DLL_MEMBER void VW_CALLING_CONV VW_MultiLinePredict(VW_HANDLE handle, VW_EXAMPLE* example_handles, size_t example_count);

#ifdef __cplusplus
}
#endif