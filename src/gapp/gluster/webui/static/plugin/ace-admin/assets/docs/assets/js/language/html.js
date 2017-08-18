/**
 * HTML patterns
 *
 * @author Craig Campbell
 * @version 1.0.9
 */
Rainbow.extend('html', [
    {
        'name': 'source.php.embedded',
        'matches': {
            2: {
                'language': 'php'
            }
        },
        'pattern': /&lt;\?=?(?!xml)(php)?([\s\S]*?)(\?&gt;)/gm
    },
    {
        'name': 'source.css.embedded',
        'matches': {
            1: {
                'matches': {
                    1: 'support.tag.style',
                    2: [
                        {
                            'name': 'entity.tag.style',
                            'pattern': /^style/g
                        },
                        {
                            'name': 'string',
                            'pattern': /('|")(.*?)(\1)/g
                        },
                        {
                            'name': 'entity.tag.style.attribute',
                            'pattern': /(\w+)/g
                        }
                    ],
                    3: 'support.tag.style'
                },
                'pattern': /(&lt;\/?)(style.*?)(&gt;)/g
            },
            2: {
                'language': 'css'
            },
            3: 'support.tag.style',
            4: 'entity.tag.style',
            5: 'support.tag.style'
        },
        'pattern': /(&lt;style.*?&gt;)([\s\S]*?)(&lt;\/)(style)(&gt;)/gm
    },
    {
        'name': 'source.js.embedded',
        'matches': {
            1: {
                'matches': {
                    1: 'support.tag.script',
                    2: [
                        {
                            'name': 'entity.tag.script',
                            'pattern': /^script/g
                        },

                        {
                            'name': 'string',
                            'pattern': /('|")(.*?)(\1)/g
                        },
                        {
                            'name': 'entity.tag.script.attribute',
                            'pattern': /(\w+)/g
                        }
                    ],
                    3: 'support.tag.script'
                },
                'pattern': /(&lt;\/?)(script.*?)(&gt;)/g
            },
            2: {
                'language': 'javascript'
            },
            3: 'support.tag.script',
            4: 'entity.tag.script',
            5: 'support.tag.script'
        },
        'pattern': /(&lt;script(?! src).*?&gt;)([\s\S]*?)(&lt;\/)(script)(&gt;)/gm
    },
    {
        'name': 'comment.html',
        'pattern': /&lt;\!--[\S\s]*?--&gt;/g
    },
    {
        'matches': {
            1: 'support.tag.open',
            2: 'support.tag.close'
        },
        'pattern': /(&lt;)|(\/?\??&gt;)/g
    },
    {
        'name': 'support.tag',
        'matches': {
            1: 'support.tag',
            2: 'support.tag.special',
            3: 'support.tag-name'
        },
        'pattern': /(&lt;\??)(\/|\!?)(\w+)/g
    },
    {
        'matches': {
            1: 'support.attribute'
        },
        'pattern': /([a-z-]+)(?=\=)/gi
    },
    {
        'matches': {
            1: 'support.operator',
            2: 'string.quote',
            3: 'string.value',
            4: 'string.quote'
        },
        'pattern': /(=)('|")(.*?)(\2)/g
    },
    {
        'matches': {
            1: 'support.operator',
            2: 'support.value'
        },
        'pattern': /(=)([a-zA-Z\-0-9]*)\b/g
    },
    {
        'matches': {
            1: 'support.attribute'
        },
        'pattern': /\s(\w+)(?=\s|&gt;)(?![\s\S]*&lt;)/g
    }
], true);
