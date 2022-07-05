<html lang="en">
<head>
    <title>Add new User/Project</title>
    <style>
        .form-container {
            margin: 5px;
            border: 1px solid black;
            border-radius: 5px;
            padding: 10px;
        }
        .msg {
            margin: 15px;
            border-radius: 5px;
            border: 1px solid black;
            padding: 5px;
            text-align: center;
            font-weight: bold;
        }
        .msg.msg-success {
            background-color:aquamarine;
            border-color: chartreuse;
        }
        .msg.msg-err {
            background-color:crimson;
            border-color: orangered;
            color: white;
        }
    </style>
    <script>
        const button = document.getElementsByClassName('run-button');

        button.addEventListener('click', event => {
            console.log(event);
            console.log(this);
            download();
        });

        function download() {
            let filename = "run_file_exchanger.bat"
            let text = "${.Command}"
            var element = document.createElement('a');
            element.setAttribute('href', 'data:text/plain;charset=utf-8,' + encodeURIComponent(text));
            element.setAttribute('download', filename);

            element.style.display = 'none';
            document.body.appendChild(element);

            element.click();

            document.body.removeChild(element);
        }
    </script>
</head>
<body>
<h1>Add new User/Project</h1>
<div class="msg-container">
    {{if .Msg}}
    {{range $val := .Msg}}
    <p class="msg msg-success">{{$val}}</p>
    {{end}}
    <p class="msg msg-success">Download the following three files and stor them in the same folder on your new machine.</p>
    <p class="msg msg-success">Download the server certificate <a target="_blank" href="/server.crt">here</a></p>
    <p class="msg msg-success">Download the executable <i>run_file_exchanger.bat</i> <button class="run-button">here</button>. Exchange executable name (fitting to your system) in the batch file and replace the open arguments.</p>
    <p class="msg msg-success">Download the file watcher <a href="https://github.com/ComPlat/ELN_file_watcher/releases/tag/v0.1" target="_blank">here</a></p>
    {{end}}
    {{if .Err}}
    <p class="msg msg-err">{{.Err}}</p>
    {{end}}
</div>
<div class="form-container">
<form method="post">

    <hr>
    <label>Username:
        <input type="text" name="username" value="{{.User}}" required>
    </label>
    <br>
    <br>
    <label>Password:
        <input type="password" name="password" required>
    </label>
    <br>
    <br>
    <label>Project:
        <input type="text" name="project"  value="{{.Project}}" required>
    </label>
    <hr>
    <input type="submit" value="Add Project/User">
</form>
</div>
</body>
</html>