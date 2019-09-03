
#include <vowpalwabbit/vwdll.h>

#ifdef __cplusplus
extern "C"
{
#endif

    typedef struct VW_ERROR {
        char* message;
    } VW_ERROR;

    VW_DLL_MEMBER VW_EXAMPLE VW_CALLING_CONV VW_ReadDSJSONExampleSafe(VW_HANDLE handle, const char* line, size_t* example_count, VW_ERROR *error);
    VW_DLL_MEMBER VW_EXAMPLE VW_CALLING_CONV VW_ReadJSONExampleSafe(VW_HANDLE handle, const char* line, size_t* example_count, VW_ERROR *error);
    VW_DLL_MEMBER void VW_CALLING_CONV VW_MultiLineLearnSafe(VW_HANDLE handle, VW_EXAMPLE* example_handles, size_t example_count, VW_ERROR *error);
    VW_DLL_MEMBER void VW_CALLING_CONV VW_MultiLinePredictSafe(VW_HANDLE handle, VW_EXAMPLE* example_handles, size_t example_count, VW_ERROR *error);
    VW_DLL_MEMBER uint32_t VW_CALLING_CONV VW_GetAction(VW_EXAMPLE e, size_t i);

#ifdef __cplusplus
}
#endif