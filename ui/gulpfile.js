const gulp = require("gulp");
const ts = require("gulp-typescript");
const tsProject = ts.createProject('tsconfig.json');

gulp.task("default", function () {
  const tsResult = gulp.src("static/ts/**/*.ts").pipe(tsProject());
  return tsResult.js.pipe(gulp.dest("static/js/"));
});
