#include <vowpalwabbit/vw.h>
#include "parse_example_json.h"
#include "lib.hpp"
#include <math.h>

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

VW_DLL_MEMBER VW_EXAMPLE VW_CALLING_CONV VW_ReadDSJSONExampleSafe(VW_HANDLE handle, const char *line, size_t *example_count, VW_ERROR *error)
{

  HANDLE_VW_ERRORS

  vw *pointer = static_cast<vw *>(handle);

  auto examples = v_init<example *>();
  examples.push_back(&VW::get_unused_example(pointer));

  DecisionServiceInteraction interaction;
  VW::read_line_decision_service_json<false>(*pointer, examples, const_cast<char *>(line), strlen(line), false,
                                             (VW::example_factory_t)&VW::get_unused_example, handle, &interaction);

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

VW_DLL_MEMBER VW_EXAMPLE VW_CALLING_CONV VW_ReadJSONExampleSafe(VW_HANDLE handle, const char *line, size_t *example_count, VW_ERROR *error)
{

  HANDLE_VW_ERRORS

  vw *pointer = static_cast<vw *>(handle);

  auto examples = v_init<example *>();
  examples.push_back(&VW::get_unused_example(pointer));

  VW::read_line_json<false>(*pointer, examples, const_cast<char *>(line), (VW::example_factory_t)&VW::get_unused_example, pointer);

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

VW_DLL_MEMBER size_t VW_CALLING_CONV VW_GetAction(VW_EXAMPLE e, size_t i)
{
  example *ex = static_cast<example *>(e);
  return ex->pred.a_s[i].action;
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

VW_DLL_MEMBER void VW_CALLING_CONV VW_EndOfPass(VW_HANDLE handle, VW_ERROR *error)
{
  HANDLE_VW_ERRORS

  vw *pointer = static_cast<vw *>(handle);
  pointer->l->end_pass();
  VW::sync_stats(*pointer);

  return;

  END_HANDLE_VW_ERRORS(error, )
}
