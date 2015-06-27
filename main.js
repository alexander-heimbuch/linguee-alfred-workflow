var search = require('./search'),
    alfred = require('alfred-workflow-nodejs'),
    actionHandler = alfred.actionHandler,
    workflow = alfred.workflow,
    Item = alfred.Item;

workflow.setName('linguee-dictionary');

(function main () {

    actionHandler.onAction('translate', function (query) {

        search(query, function (results) {
            results.forEach(function (result) {
                if (result.translation === '') {
                    workflow.addItem(new Item({
                        title: result.suggestion,
                        subtitle: '',
                        valid: true,
                        arg: result.href
                    }));
                } else {
                    workflow.addItem(new Item({
                        title: result.translation,
                        subtitle: result.suggestion,
                        valid: true,
                        arg: result.href
                    }))
                }
            });
            workflow.feedback();
        });
    });

    alfred.run();
})();
