
#include <stdint.h>
#include <vowpalwabbit/vwdll.h>

#ifdef __cplusplus
extern "C"
{
#endif

    typedef void *VW_EXAMPLE_POOL_HANDLE;

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

    VW_DLL_PUBLIC VW_EXAMPLE VW_CALLING_CONV VW_ReadDSJSONExampleSafe(VW_HANDLE handle, VW_EXAMPLE_POOL_HANDLE example_pool, const char *line, size_t *example_count, VW_ERROR *error);
    VW_DLL_PUBLIC VW_EXAMPLE VW_CALLING_CONV VW_ReadJSONExampleSafe(VW_HANDLE handle, VW_EXAMPLE_POOL_HANDLE example_pool, const char *line, size_t *example_count, VW_ERROR *error);

    VW_DLL_PUBLIC void VW_CALLING_CONV VW_MultiLineLearnSafe(VW_HANDLE handle, VW_EXAMPLE *example_handles, size_t example_count, VW_ERROR *error);
    VW_DLL_PUBLIC void VW_CALLING_CONV VW_MultiLinePredictSafe(VW_HANDLE handle, VW_EXAMPLE *example_handles, size_t example_count, VW_ERROR *error);

    VW_DLL_PUBLIC size_t VW_CALLING_CONV VW_GetAction(VW_EXAMPLE e, size_t i);
    VW_DLL_PUBLIC float VW_CALLING_CONV VW_GetLoss(VW_EXAMPLE e);
    VW_DLL_PUBLIC size_t VW_CALLING_CONV VW_GetCBCostLength(VW_EXAMPLE e);
    VW_DLL_PUBLIC float VW_CALLING_CONV VW_GetCBCost(VW_EXAMPLE e, size_t i);
    VW_DLL_PUBLIC size_t VW_CALLING_CONV VW_GetMultiClassPrediction(VW_EXAMPLE e);
    VW_DLL_PUBLIC size_t VW_CALLING_CONV VW_GetScalarLength(VW_EXAMPLE e);
    VW_DLL_PUBLIC float VW_CALLING_CONV VW_GetScalar(VW_EXAMPLE e, size_t i);

    VW_DLL_PUBLIC float VW_CALLING_CONV VW_GetLearningRate(VW_HANDLE handle, VW_ERROR *error);
    VW_DLL_PUBLIC void VW_CALLING_CONV VW_SyncStats(VW_HANDLE handle, VW_ERROR *error);
    VW_DLL_PUBLIC VW_PERFORMANCE_STATS VW_CALLING_CONV VW_PerformanceStats(VW_HANDLE handle, VW_ERROR *error);
    VW_DLL_PUBLIC void VW_CALLING_CONV VW_EndOfPass(VW_HANDLE handle, VW_ERROR *error);

    VW_DLL_PUBLIC VW_EXAMPLE_POOL_HANDLE VW_CALLING_CONV VW_CreateExamplePool(VW_HANDLE handle);
    VW_DLL_PUBLIC void VW_CALLING_CONV VW_ReleaseExamplePool(VW_EXAMPLE_POOL_HANDLE handle);
    VW_DLL_PUBLIC void VW_CALLING_CONV VW_ReturnExampleToPool(VW_EXAMPLE_POOL_HANDLE pool_handle, VW_EXAMPLE e);

#ifdef __cplusplus
}
#endif