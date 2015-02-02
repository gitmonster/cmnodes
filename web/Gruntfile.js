module.exports = function(grunt) {

    grunt.initConfig({
        'connect': {
            server: {
                options: {
                    open: true,
                    keepalive: true,
                    port: 8000,
                    base: '.'
                }
            }
        },
        'watch': {
            options: {
                livereload: true,
            },
            html: {
                files: ['index.html', 'src/**/*', '!dist'],
                tasks: ['build'],
            },
        },
        'gh-pages': {
            options: {
                clone: 'bower_components/cmnodes'
            },
            src: [
                'bower_components/**/*',
                '!bower_components/cmnodes/**/*',
                'demo/*', 'src/*', 'index.html'
            ]
        },
        'copy': {
            main: {
                files: [{
                    expand: true,
                    cwd: 'bower_components',
                    src: ['**/*'],
                    dest: 'dist/vendor/'
                }],
            },
        },
        'replace': {
            index: {
                src: ['index.html'],
                dest: 'dist/',
                replacements: [{
                    from: 'bower_components',
                    to: 'vendor'
                }, {
                    from: /href=\"src/g,
                    to: 'href=\"modules'
                }]
            },
            modules: {
                src: ['src/**/*'],
                dest: 'dist/modules/',
                replacements: [{
                    from: 'bower_components',
                    to: 'vendor'
                }]
            }
        }
    });

    grunt.loadNpmTasks('grunt-contrib-copy');
    grunt.loadNpmTasks('grunt-contrib-connect');
    grunt.loadNpmTasks('grunt-gh-pages');
    grunt.loadNpmTasks('grunt-text-replace');
    grunt.loadNpmTasks('grunt-contrib-watch');

    grunt.registerTask('build', ['replace', 'copy']);
    grunt.registerTask('deploy', ['gh-pages']);
    grunt.registerTask('default', ['build', 'connect', 'watch']);

    grunt.event.on('watch', function(action, filepath, target) {
        grunt.log.writeln(target + ': ' + filepath + ' has ' + action);
    });

};