
#include <stdint.h>
#include <vowpalwabbit/vwdll.h>

#ifdef __cplusplus
extern "C"
{
#endif

    typedef struct VW_ERROR
    {
        char *message;
    } VW_ERROR;

    typedef struct VW_PERFORMANCE_STATS
    {
        uint64_t current_pass;
        uint64_t number_of_features;
        uint64_t number_of_examples;
        double weighted_example_sum;
        double weighted_label_sum;
        double average_loss;
        double best_constant;
        double best_constant_loss;
    } VW_PERFORMANCE_STATS;

    VW_DLL_MEMBER VW_EXAMPLE VW_CALLING_CONV VW_ReadDSJSONExampleSafe(VW_HANDLE handle, const char *line, size_t *example_count, VW_ERROR *error);
    VW_DLL_MEMBER VW_EXAMPLE VW_CALLING_CONV VW_ReadJSONExampleSafe(VW_HANDLE handle, const char *line, size_t *example_count, VW_ERROR *error);

    VW_DLL_MEMBER void VW_CALLING_CONV VW_MultiLineLearnSafe(VW_HANDLE handle, VW_EXAMPLE *example_handles, size_t example_count, VW_ERROR *error);
    VW_DLL_MEMBER void VW_CALLING_CONV VW_MultiLinePredictSafe(VW_HANDLE handle, VW_EXAMPLE *example_handles, size_t example_count, VW_ERROR *error);

    VW_DLL_MEMBER size_t VW_CALLING_CONV VW_GetAction(VW_EXAMPLE e, size_t i);

    VW_DLL_MEMBER VW_PERFORMANCE_STATS VW_CALLING_CONV VW_PerformanceStats(VW_HANDLE handle, VW_ERROR *error);
    VW_DLL_MEMBER void VW_CALLING_CONV VW_EndOfPass(VW_HANDLE handle, VW_ERROR *error);

#ifdef __cplusplus
}
#endif