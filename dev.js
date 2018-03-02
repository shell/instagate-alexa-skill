// Script to generate UTTERANCES file out of intents.json
var fs = require('fs')
var intents = JSON.parse(fs.readFileSync('intents.json'))
var text = '';
intents['languageModel']['intents'].forEach((intent) => {
    if (!intent.name.match('AMAZON')) {
        text += intent.samples.map(sample => `${intent.name} ${sample}`).join("\n") + "\n\n";
    }
});
fs.writeFileSync('UTTERANCES', text.replace(/\n+$/, ''));