#include <vowpalwabbit/vw.h>
#include <vowpalwabbit/parse_example_json.h>
#include "lib.hpp"
#include <math.h>
#include <iostream>

#define HANDLE_VW_ERRORS \
  try                    \
  {
#define END_HANDLE_VW_ERRORS(errVar, retVal)  \
  }                                           \
  catch (VW::vw_exception const &e)           \
  {                                           \
    auto msg = e.what();                      \
    auto err = VW_ERROR{                      \
        .message = new char[strlen(msg) + 1], \
    };                                        \
    std::strcpy(err.message, msg);            \
    *errVar = err;                            \
    return retVal;                            \
  }                                           \
  catch (const std::exception &e)             \
  {                                           \
    auto msg = e.what();                      \
    auto err = VW_ERROR{                      \
        .message = new char[strlen(msg) + 1], \
    };                                        \
    std::strcpy(err.message, msg);            \
    *errVar = err;                            \
    return retVal;                            \
  }

struct ExamplePool
{
  vw *_vw;
  std::vector<example *> _example_pool;

  example *get_or_create_example()
  {
    if (_example_pool.size() == 0)
    {
      auto ex = VW::alloc_examples(0, 1);
      _vw->p->lp.default_label(&ex->l);

      return ex;
    }

    // get last element
    example *ex = _example_pool.back();
    _example_pool.pop_back();

    VW::empty_example(*_vw, *ex);
    _vw->p->lp.default_label(&ex->l);

    return ex;
  }
};

example &get_or_create_example_f(VW_EXAMPLE_POOL_HANDLE example_pool_handle)
{
  ExamplePool *pool = static_cast<ExamplePool *>(example_pool_handle);
  return *pool->get_or_create_example();
}

VW_DLL_MEMBER VW_EXAMPLE VW_CALLING_CONV VW_ReadDSJSONExampleSafe(VW_HANDLE handle, VW_EXAMPLE_POOL_HANDLE example_pool, const char *line, size_t *example_count, VW_ERROR *error)
{

  HANDLE_VW_ERRORS

  ExamplePool *pool = static_cast<ExamplePool *>(example_pool);
  vw *pointer = static_cast<vw *>(handle);

  auto examples = v_init<example *>();
  examples.push_back(pool->get_or_create_example());

  DecisionServiceInteraction interaction;
  VW::read_line_decision_service_json<false>(*pointer, examples, const_cast<char *>(line), strlen(line), false,
                                             (VW::example_factory_t)get_or_create_example_f, example_pool, &interaction);

  VW::setup_examples(*pointer, examples);

  *example_count = examples.size();

  auto exmpl = (VW_EXAMPLE *)malloc(sizeof(VW_EXAMPLE) * examples.size());
  for (size_t i = 0; i < examples.size(); ++i)
  {
    exmpl[i] = examples[i];
  }
  examples.delete_v();
  return exmpl;

  END_HANDLE_VW_ERRORS(error, NULL)
}

VW_DLL_MEMBER VW_EXAMPLE VW_CALLING_CONV VW_ReadJSONExampleSafe(VW_HANDLE handle, VW_EXAMPLE_POOL_HANDLE example_pool, const char *line, size_t *example_count, VW_ERROR *error)
{

  HANDLE_VW_ERRORS

  ExamplePool *pool = static_cast<ExamplePool *>(example_pool);
  vw *pointer = static_cast<vw *>(handle);

  auto examples = v_init<example *>();
  examples.push_back(pool->get_or_create_example());

  VW::read_line_json<false>(*pointer, examples, const_cast<char *>(line), (VW::example_factory_t)get_or_create_example_f, example_pool);

  VW::setup_examples(*pointer, examples);

  *example_count = examples.size();

  auto exmpl = (VW_EXAMPLE *)malloc(sizeof(VW_EXAMPLE) * examples.size());
  for (size_t i = 0; i < examples.size(); ++i)
  {
    exmpl[i] = examples[i];
  }
  examples.delete_v();
  return exmpl;

  END_HANDLE_VW_ERRORS(error, NULL)
}

VW_DLL_MEMBER void VW_CALLING_CONV VW_MultiLineLearnSafe(VW_HANDLE handle, VW_EXAMPLE *example_handles, size_t example_count, VW_ERROR *error)
{

  HANDLE_VW_ERRORS

  vw *pointer = static_cast<vw *>(handle);
  auto examples = (example **)example_handles;

  multi_ex examples_vector(examples, examples + example_count);

  pointer->learn(examples_vector);

  return;

  END_HANDLE_VW_ERRORS(error, )
}

VW_DLL_MEMBER void VW_CALLING_CONV VW_MultiLinePredictSafe(VW_HANDLE handle, VW_EXAMPLE *example_handles, size_t example_count, VW_ERROR *error)
{

  HANDLE_VW_ERRORS

  vw *pointer = static_cast<vw *>(handle);
  auto examples = (example **)example_handles;

  multi_ex examples_vector(examples, examples + example_count);

  pointer->predict(examples_vector);

  return;

  END_HANDLE_VW_ERRORS(error, )
}

VW_DLL_MEMBER size_t VW_CALLING_CONV VW_GetScalarLength(VW_EXAMPLE e)
{
  example *ex = static_cast<example *>(e);
  return ex->pred.scalars.size();
}

VW_DLL_MEMBER float VW_CALLING_CONV VW_GetScalar(VW_EXAMPLE e, size_t i)
{
  example *ex = static_cast<example *>(e);
  return ex->pred.scalars[i];
}

VW_DLL_MEMBER size_t VW_CALLING_CONV VW_GetAction(VW_EXAMPLE e, size_t i)
{
  example *ex = static_cast<example *>(e);
  return ex->pred.a_s[i].action;
}

VW_DLL_MEMBER float VW_CALLING_CONV VW_GetCBCost(VW_EXAMPLE e, size_t i)
{
  example *ex = static_cast<example *>(e);
  return ex->l.cb.costs[i].cost;
}

VW_DLL_MEMBER size_t VW_CALLING_CONV VW_GetCBCostLength(VW_EXAMPLE e)
{
  example *ex = static_cast<example *>(e);
  return ex->l.cb.costs.size();
}

VW_DLL_MEMBER size_t VW_CALLING_CONV VW_GetMultiClassPrediction(VW_EXAMPLE e)
{
  example *ex = static_cast<example *>(e);
  return ex->pred.multiclass;
}

VW_DLL_MEMBER float VW_CALLING_CONV VW_GetLoss(VW_EXAMPLE e)
{
  example *ex = static_cast<example *>(e);
  return ex->loss;
}

VW_DLL_MEMBER VW_PERFORMANCE_STATS VW_CALLING_CONV VW_PerformanceStats(VW_HANDLE handle, VW_ERROR *error)
{
  HANDLE_VW_ERRORS

  vw *pointer = static_cast<vw *>(handle);

  auto stats = VW_PERFORMANCE_STATS{};

  stats.current_pass = pointer->current_pass;
  stats.number_of_examples = pointer->sd->example_number;
  stats.weighted_example_sum = pointer->sd->weighted_examples();
  stats.weighted_label_sum = pointer->sd->weighted_labels;

  if (pointer->holdout_set_off)
    if (pointer->sd->weighted_labeled_examples > 0)
      stats.average_loss = pointer->sd->sum_loss / pointer->sd->weighted_labeled_examples;
    else
      stats.average_loss = NAN;
  else if ((pointer->sd->holdout_best_loss == FLT_MAX) || (pointer->sd->holdout_best_loss == FLT_MAX * 0.5))
    stats.average_loss = NAN;
  else
    stats.average_loss = pointer->sd->holdout_best_loss;

  float best_constant;
  float best_constant_loss;
  if (get_best_constant(*pointer, best_constant, best_constant_loss))
  {
    stats.best_constant = best_constant;
    if (best_constant_loss != FLT_MIN)
    {
      stats.best_constant_loss = best_constant_loss;
    }
  }

  stats.number_of_features = pointer->sd->total_features;

  return stats;

  END_HANDLE_VW_ERRORS(error, VW_PERFORMANCE_STATS{})
}

VW_DLL_MEMBER void VW_CALLING_CONV VW_SyncStats(VW_HANDLE handle, VW_ERROR *error)
{
  HANDLE_VW_ERRORS

  vw *pointer = static_cast<vw *>(handle);

  VW::sync_stats(*pointer);

  return;

  END_HANDLE_VW_ERRORS(error, )
}

VW_DLL_MEMBER void VW_CALLING_CONV VW_EndOfPass(VW_HANDLE handle, VW_ERROR *error)
{
  HANDLE_VW_ERRORS

  vw *pointer = static_cast<vw *>(handle);

  pointer->do_reset_source = false;
  pointer->passes_complete++;

  pointer->current_pass++;
  pointer->l->end_pass();

  VW::sync_stats(*pointer);

  return;

  END_HANDLE_VW_ERRORS(error, )
}

VW_DLL_MEMBER VW_EXAMPLE_POOL_HANDLE VW_CALLING_CONV VW_CreateExamplePool(VW_HANDLE handle)
{
  vw *pointer = static_cast<vw *>(handle);
  auto pool = new ExamplePool;

  pool->_vw = pointer;

  return pool;
}

VW_DLL_MEMBER void VW_CALLING_CONV VW_ReleaseExamplePool(VW_EXAMPLE_POOL_HANDLE handle)
{

  ExamplePool *pool = static_cast<ExamplePool *>(handle);
  for (auto &&ex : pool->_example_pool)
  {
    VW::dealloc_example(pool->_vw->p->lp.delete_label, *ex);
    ::free_it(ex);
  }

  pool->_example_pool.empty();
  pool->_vw = NULL;

  delete pool;
}

VW_DLL_MEMBER void VW_CALLING_CONV VW_ReturnExampleToPool(VW_EXAMPLE_POOL_HANDLE pool_handle, VW_EXAMPLE e)
{
  ExamplePool *pool = static_cast<ExamplePool *>(pool_handle);
  example *ex = static_cast<example *>(e);

  pool->_example_pool.emplace_back(ex);
}