#include <vowpalwabbit/vw.h>
#include "parse_example_json.h"
#include "lib.hpp"


#define HANDLE_VW_ERRORS                                           \
  try {
#define END_HANDLE_VW_ERRORS(errVar, retVal)                       \
  }                                                                \
  catch (VW::vw_exception const& e) {                              \
    auto msg = e.what();                                           \
    auto err = VW_ERROR{                                           \
        .message = new char[strlen(msg)+1],                        \
    };                                                             \
    std::strcpy(err.message, msg);                                 \
    *errVar = err;                                                 \
    return retVal;                                                 \
  }                                                                \
  catch (const std::exception& e) {                                \
    auto msg = e.what();                                           \
    auto err = VW_ERROR{                                           \
        .message = new char[strlen(msg)+1],                        \
    };                                                             \
    std::strcpy(err.message, msg);                                 \
    *errVar = err;                                                 \
    return retVal;                                                 \
  }

VW_DLL_MEMBER VW_EXAMPLE VW_CALLING_CONV VW_ReadDSJSONExampleSafe(VW_HANDLE handle, const char* line, size_t* example_count, VW_ERROR *error) {
    
    HANDLE_VW_ERRORS

    vw * pointer = static_cast<vw*>(handle);
    
    auto examples = v_init<example*>();
    examples.push_back(&VW::get_unused_example(pointer));

    DecisionServiceInteraction interaction;
    VW::read_line_decision_service_json<false>(*pointer, examples, const_cast<char*>(line), strlen(line), false,
        (VW::example_factory_t)&VW::get_unused_example, handle, &interaction);


    VW::setup_examples(*pointer, examples);

    *example_count = examples.size();

    auto exmpl = (VW_EXAMPLE*)malloc(sizeof(VW_EXAMPLE) * examples.size());
    for (size_t i = 0; i < examples.size(); ++i) {
        exmpl[i] = examples[i];
    }
    examples.delete_v();
    return exmpl;

    END_HANDLE_VW_ERRORS(error, NULL)
}


VW_DLL_MEMBER VW_EXAMPLE VW_CALLING_CONV VW_ReadJSONExampleSafe(VW_HANDLE handle, const char* line, size_t* example_count, VW_ERROR *error) {
    
    HANDLE_VW_ERRORS
    
    vw * pointer = static_cast<vw*>(handle);

    auto examples = v_init<example*>();
    examples.push_back(&VW::get_unused_example(pointer));

    VW::read_line_json<false>(*pointer, examples, const_cast<char*>(line), (VW::example_factory_t)&VW::get_unused_example, pointer);

    VW::setup_examples(*pointer, examples);

    *example_count = examples.size();

    auto exmpl = (VW_EXAMPLE*)malloc(sizeof(VW_EXAMPLE) * examples.size());
    for (size_t i = 0; i < examples.size(); ++i) {
        exmpl[i] = examples[i];
    }
    examples.delete_v();
    return exmpl;

    END_HANDLE_VW_ERRORS(error, NULL)
}


VW_DLL_MEMBER void VW_CALLING_CONV VW_MultiLineLearnSafe(VW_HANDLE handle, VW_EXAMPLE* example_handles, size_t example_count, VW_ERROR *error) {
    
    HANDLE_VW_ERRORS
    
    vw * pointer = static_cast<vw*>(handle);
    auto examples = (example**)example_handles;

    multi_ex examples_vector(examples, examples + example_count);

    pointer->learn(examples_vector);

    return;

    END_HANDLE_VW_ERRORS(error,)
}


VW_DLL_MEMBER void VW_CALLING_CONV VW_MultiLinePredictSafe(VW_HANDLE handle, VW_EXAMPLE* example_handles, size_t example_count, VW_ERROR *error) {
    
    HANDLE_VW_ERRORS
    
    vw* pointer = static_cast<vw*>(handle);
    auto examples = (example**)example_handles;

    multi_ex examples_vector(examples, examples + example_count);

    pointer->predict(examples_vector);

    return;

    END_HANDLE_VW_ERRORS(error,)
}

VW_DLL_MEMBER size_t VW_CALLING_CONV VW_GetAction(VW_EXAMPLE e, size_t i) { 
  example * ex = static_cast<example*>(e);
  return ex->pred.a_s[i].action;
}
