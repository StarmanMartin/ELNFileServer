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
            background-color: aquamarine;
            border-color: chartreuse;
        }

        .msg.msg-err {
            background-color: crimson;
            border-color: orangered;
            color: white;
        }
    </style>
    <script>
        function download() {
            let filename = "run_file_exchanger.bat"
            let text = "{{.Command}}"
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
    <div class="msg msg-success">
        <p>Download the following</p>
        <ol>
            <li>run_file_exchanger.bat</li>
            <li>server.crt</li>
            <li>efw_[architecture].exe .crt</li>
        </ol>
        <p>and store them in the same folder on your new machine.</p>
        <p>Download the executable <i>run_file_exchanger.bat</i>
            <button onclick="download()">here</button>
            . Exchange executable name (fitting to your system) in the batch file and replace the open arguments.
        </p>
    </div>
    {{end}}
    {{if .Err}}
    <p class="msg msg-err">{{.Err}}</p>
    {{end}}
    <div class="msg">
        <p>Download the server certificate <a target="_blank" href="/server.crt">here</a></p>
        <p>Download the file watcher <a href="https://github.com/ComPlat/ELN_file_watcher/releases/tag/v0.1"
                                        target="_blank">here</a></p>
    </div>
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
            <input type="text" name="project" value="{{.Project}}" required>
        </label>
        <hr>
        <input type="submit" value="Add Project/User">
    </form>
</div>
</body>
</html>