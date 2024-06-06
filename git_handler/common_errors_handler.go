package githandler

import (
	"fmt"
	"strings"

	// "sync"

	"log"

	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
	filehandler "github.com/it-shiloheye/ftp_system_lib/file_handler"
)

func handle_common_git_errors(ctx ftp_context.Context, directory string, stderr string, cmd_err error) (retry bool, err *ftp_context.LogItem) {
	loc := "handle_common_git_errors"
	var buf []byte
	fmt.Println(loc)
	dec_child_count_f(ctx)
	defer dec_child_count_f(ctx)

	if strings.Contains(stderr, "not a git repository") {
		log.Println("not a git repository")
		c := strings_split("git init .", " ")
		buf, stderr, err = ExecuteCommand(ctx, directory, c[0], c[1:]...)

		if err != nil {
			log.Println(err)

			return handle_common_git_errors(ctx, directory, stderr, err.AppendParentError(err))
		}
		log.Println(string(buf))

		Fo := filehandler.NewFileBasic(directory + "/.gitignore").Open()
		if err = Fo.Err; err != nil {
			return
		}
		fo_2 := filehandler.NewFileBasic("./data/templates/.gitignore").Open()
		if err = fo_2.Err; err != nil {
			return
		}
		buf, _ = fo_2.ReadAll()
		if fo_2.Err != nil {
			return
		}

		Fo.Write(buf)
		if Fo.Err != nil {
			return false, Fo.Err
		}

		return true, nil
	}
	if stderr == "" {
		log.Println("No real error")
		// c := strings_split("rm "+directory+"/.git/index.lock", " ")
		// buf, stderr, err = ExecuteCommand(ctx, directory, c[0], c[1:]...)

		// if err != nil {
		// 	log.Println(err)
		// 	return handle_common_git_errors(ctx, directory, stderr, ftp_context.NewLogItem(loc,true).AppendParentError(err))
		// }
		// log.Println(string(buf))

		return false, nil
	}

	if strings.Contains(stderr, "Another git process seems to be running in this repository") {
		log.Println("Another git process seems to be running")
		c := strings_split("taskkill -im git -f", " ")
		buf, stderr, err = ExecuteCommand(ctx, directory, c[0], c[1:]...)

		if err != nil {
			log.Println(err)
			return handle_common_git_errors(ctx, directory, stderr, ftp_context.NewLogItem(loc, true).Set("original_error", cmd_err).AppendParentError(err))
		}
		log.Println(string(buf))
		retry = true
		return
	}

	if strings.Contains(stderr, "Another git process seems to be running in this repository") {
		log.Println("Another git process seems to be running")
		c := strings_split("taskkill -im git -f", " ")
		buf, stderr, err = ExecuteCommand(ctx, directory, c[0], c[1:]...)

		if err != nil {
			log.Println(err)
			return handle_common_git_errors(ctx, directory, stderr, ftp_context.NewLogItem(loc, true).Set("original_error", cmd_err).AppendParentError(err))
		}
		log.Println(string(buf))
		retry = true
		return
	}
	return
}
