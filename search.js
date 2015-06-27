var $ = require('cheerio'),
    request = require('request'),
    iconv = require('iconv'),
    entities = require('entities');

module.exports = function (query, cb) {
    request('http://www.linguee.de/deutsch-englisch/search?qe=' + encodeURIComponent(query) + '&source=auto', {'encoding': 'binary'}, function (err, resp, html) {
        var $response,
            ic = iconv.Iconv('iso-8859-15', 'utf-8'),
            results = [];

        if (err) return cb(results);

        html = ic.convert(new Buffer(html, 'binary')).toString('utf-8');

        $response = $.load(html);
        $response('.autocompletion_item').map(function (i, suggestion) {
            var $suggestion = $(suggestion),
                suggestion,
                translations = [];

            $suggestion.find('.translation_item').map(function (i, translation) {
                var $translation = $(translation);

                $translation.find('.wordtype').remove();
                $translation.find('.sep').remove();

                translations.push(entities.decodeXML(
                    $translation
                        .text()
                        .replace(/(\r\n|\n|\r)/gm, '')
                        .trim()
                ));
            });

            suggestion = entities.decodeXML($suggestion.find('.main_item').text());

            results.push({
                'suggestion': suggestion,
                'translation': translations.join(', '),
                'href': 'http://www.linguee.de/deutsch-englisch/search?source=auto&query=' + encodeURIComponent(suggestion)
            });
        });

        cb(results);
    });
}
