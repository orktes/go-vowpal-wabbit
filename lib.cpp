#include <vowpalwabbit/vw.h>
#include "parse_example_json.h"
#include "lib.hpp"
#include <iostream>

VW_DLL_MEMBER VW_EXAMPLE VW_CALLING_CONV VW_ReadDSJSONExample(VW_HANDLE handle, const char* line, size_t* example_count) {
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
}


VW_DLL_MEMBER VW_EXAMPLE VW_CALLING_CONV VW_ReadJSONExample(VW_HANDLE handle, const char* line, size_t* example_count) {
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
}


VW_DLL_MEMBER void VW_CALLING_CONV VW_MultiLineLearn(VW_HANDLE handle, VW_EXAMPLE* example_handles, size_t example_count) {
    vw * pointer = static_cast<vw*>(handle);
    
    auto examples = (example**)example_handles;

    multi_ex examples_vector(examples, examples + example_count);

    pointer->learn(examples_vector);

    return;
}


VW_DLL_MEMBER void VW_CALLING_CONV VW_MultiLinePredict(VW_HANDLE handle, VW_EXAMPLE* example_handles, size_t example_count) {
    vw* pointer = static_cast<vw*>(handle);
    auto examples = (example**)example_handles;

    multi_ex examples_vector(examples, examples + example_count);

    pointer->predict(examples_vector);

    return;
}

