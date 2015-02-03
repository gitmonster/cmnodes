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
                clone: 'lib/cmnodes'
            },
            src: [
                'lib/**/*',
                '!lib/cmnodes/**/*',
                'demo/*', 'src/*', 'index.html'
            ]
        },
        'copy': {
            lib: {
                files: [{
                    expand: true,
                    cwd: 'lib',
                    src: ['**/*'],
                    dest: 'dist/lib/'
                }],
            },
            modules: {
                files: [{
                    expand: true,
                    cwd: 'src',
                    src: ['**/*'],
                    dest: 'dist/src/'
                }],
            },
            index: {
                files: [{
                    src: ['index.html'],
                    dest: 'dist/'
                }],
            },
        }
    });
    
    grunt.loadNpmTasks('grunt-contrib-copy');
    grunt.loadNpmTasks('grunt-contrib-connect');
    grunt.loadNpmTasks('grunt-gh-pages');
    grunt.loadNpmTasks('grunt-contrib-watch');

    grunt.registerTask('build', ['copy']);
    grunt.registerTask('deploy', ['gh-pages']);
    grunt.registerTask('default', ['build', 'connect', 'watch']);

    grunt.event.on('watch', function(action, filepath, target) {
        grunt.log.writeln(target + ': ' + filepath + ' has ' + action);
    });

};