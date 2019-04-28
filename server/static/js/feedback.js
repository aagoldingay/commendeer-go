function submitFeedback() {
    // configure data objects
    var form = document.getElementById("questionnaire")
    form.method = "POST";
    form.action = "/thanks";
    //form.children[0] = id
    //form.children[1] = accesscode
    var questions = form.querySelectorAll(".question");
    var query = [];

    // populate query from questions
    for (var i = 0; i < questions.length; i++) {
        // questions.children :=
            // 0 = hidden qid
            // 1 = hidden qtype
            // 2 = p
            // 3 = answer / options
        var question = {
            questionid: parseInt(questions[i].children[0].value, 10),
            type: parseInt(questions[i].children[1].value, 10),
            answer: "",
            options: []
        };
        
        if (question.type < 3) { // if question = multichoice
            var opts = questions[i].children[3].children;
            for (var j = 0; j < opts.length; j+=3) {
                if (opts[j].checked) {
                    var optval = parseInt(opts[j].value, 10);
                    question.options.push(optval)
                }
            }
        } else {
            question.answer = questions[i].children[3].children[0].value
        }
        query.push(question)
    }
    var questionnaire = {
        id: parseInt(form.children[0].value),
        accesscode: form.children[1].value,
        questions: query
    };
    console.log(query)
    console.log(questionnaire)
    console.log(JSON.stringify(questionnaire))

    // set input "request" with string json
    document.getElementById("request").value = JSON.stringify(questionnaire)
    form.submit();
}