current_dir = get_dir(current_filepath())

# arch_install("yay")
# arch_install("powerpill")

# def install_packages():
#     darwin_install(name="alacritty", cask=True)
#     darwin_install(name="clojure")
#     darwin_install(name="go")
#     darwin_install(name="tmux")
#     darwin_install(name="vim")
#     darwin_install(name="bat")
#     darwin_install(name="bitwarden-cli")
#     darwin_install(name="dive")
#     darwin_install(name="direnv")
#     darwin_install(name="exa")
#     darwin_install(name="fd")
#     darwin_install(name="fzf")
#     darwin_install(name="ghq")
#     darwin_install(name="jq")
#     darwin_install(name="make")
#     darwin_install(name="nkf")
#     darwin_install(name="ripgrep")
#     darwin_install(name="procs")
#     darwin_install(name="kubectx")
#     darwin_install(name="starship")
#     darwin_install(name="less")
#     darwin_install(name="reattach-to-user-namespace")
#     darwin_install(name="stern")
#     darwin_install(name="wget")
#     darwin_install(name="alfred", cask=True)
#     darwin_install(name="android-studio", cask=True)
#     darwin_install(name="bartender", cask=True)
#     darwin_install(name="clipy", cask=True)
#     darwin_install(name="dash", cask=True)
#     darwin_install(name="docker", cask=True)
#     darwin_install(name="homebrew/cask-fonts/font-cica", cask=True)
#     darwin_install(name="homebrew/cask-fonts/font-myrica", cask=True)
#     darwin_install(name="google-chrome", cask=True)
#     darwin_install(name="google-cloud-sdk", cask=True)
#     darwin_install(name="hyperswitch", cask=True)
#     darwin_install(name="inkdrop", cask=True)
#     darwin_install(name="mplayerx", cask=True)
#     darwin_install(name="quicklook-csv", cask=True)
#     darwin_install(name="qlmarkdown", cask=True)
#     darwin_install(name="slack", cask=True)
#     darwin_install(name="the-unarchiver", cask=True)
#     darwin_install(name="visual-studio-code", cask=True)
#     darwin_install(name="git")
#     darwin_install(name="tig")
#     darwin_install(name="hammerspoon", cask=True)
#     darwin_install(name="fish")
#     darwin_install(name="docker-edge", cask=True)
#     arch_multiple_install(names=[
#         "autoconf",
#         "automake",
#         "base-devel",
#         "binutils",
#         "fakeroot",
#         "python-pynvim",
#         "gcc",
#         "lm_sensors",
#         "libinput-gestures",
#         "google-chrome",
#         "station",
#         "gnome-tweak-tool",
#         "bitwarden",
#         "android-studio",
#         "noto-fonts",
#         "noto-fonts-cjk",
#         "noto-fonts-extra",
#         "ttf-cica",
#         "ttf-hackgen",
#         "patch",
#         "kubectl",
#         "stern-bin",
#         "minikube",
#         "skaffold",
#         "fcitx",
#         "fcitx-configtool",
#         "fcitx-mozc",
#         "fcitx-gtk3",
#         "materia-gtk-theme",
#         "gnome-themes-extra",
#         "gtk-engine-murrine",
#         "arc-gtk-theme",
#         "bat",
#         "bitwarden-cli",
#         "dive",
#         "direnv",
#         "exa",
#         "fd",
#         "fzf",
#         "ghq",
#         "jq",
#         "make",
#         "nkf",
#         "ripgrep",
#         "procs",
#         "kubectx",
#         "starship",
#         "git",
#         "alacritty",
#         "go",
#         "tmux",
#         "gvim", # +clipboardのため
#         "clojure",
#         "tig",
#         "glibc",
#         "libcanberra",
#         "gvfs",
#         "visual-studio-code-bin"
#     ])
#     arch_uninstall("gnome-terminal")
#     arch_uninstall("firefox")
#     arch_uninstall("gtkhash-nautilus")
#     arch_uninstall("gtkhash")
#     arch_uninstall("hexchat")
#     arch_uninstall("gnome-calculator")
#     arch_install("docker")
#     arch_install("docker-compose")
#     arch_install("docker-machine")
#     arch_install("libvirt")
#     arch_install("qemu")
#     arch_install("ebtables")
#     arch_install("dnsmasq")
#     arch_install("docker-machine-driver-kvm2")
#     arch_install(name="fish")
# 
# 
# def arch_base_setup():
#     if is_arch_linux():
#         # command("LANG=C xdg-user-dirs-gtk-update")
#         # command("pacman-mirrors --fasttrack")
#         command("pacman", ["-Syu", "--noconfirm"])
#         arch_install(name="yay")
#         create_directory("~/.config")
#         copy_file("~/.xprofile", join_filepath(current_dir, "config", "xprofile"))
# 
# def darwin_base_setup():
#     if is_darwin():
#         darwin_defaults = [
#             dict(domain="NSGlobalDomain", key="NSTextShowsControlCharacters", type="bool", value="true"), # ASCII制御文字を表示する
#             dict(domain="NSGlobalDomain", key="NSQuitAlwaysKeepsWindows", type="bool", value="false"), # アプリケーションを終了して再度開くときにウィンドウを復元しない
#             dict(domain="NSGlobalDomain", key="NSDisableAutomaticTermination", type="bool", value="false"), # automatic termination機能の無効化
#             dict(domain="NSGlobalDomain", key="com.apple.swipescrolldirection", type="bool", value="true"), # スクロールの方向 ナチュラル
#             dict(domain="NSGlobalDomain", key="KeyRepeat", type="int", value="2"), # キーリピートの速さ 最速
#             dict(domain="NSGlobalDomain", key="InitialKeyRepeat", type="int", value="15"),
#             dict(domain="NSGlobalDomain", key="AppleShowAllExtensions", type="bool", value="true"), # すべてのファイルの拡張子を表示
#             dict(domain="com.apple.BazelServices", key="kDim", type="bool", value="true"), # 環境光が暗い場合にキーボードの輝度を調整
#             dict(domain="com.apple.BazelServices", key="kDimTime", type="int", value="15"), # 発行した状態で待機する時間 15秒
#             dict(domain="com.apple.screensaver", key="askForPassword", type="int", value="1"), # スクリーンセーバー解除時にパスワードを要求する
#             dict(domain="com.apple.screensaver", key="askForPasswordDelay", type="float", value="0"), # スクリーンセーバーに入ってから何秒後からパスワードを要求するか 0秒
#             dict(domain="com.apple.screencapture", key="location", type="string", value="$HOME/Desktop"), # スクリーンショットの出力先
#             dict(domain="com.apple.screencapture", key="disable-shadow", type="bool", value="true"), # スクリーンショットの影付き効果なし
#             dict(domain="com.apple.dashboard", key="mcx-disabled", type="bool", value="true"), # Dashboardを無効にする
#             dict(domain="com.apple.dock", key="mru-spaces", type="bool", value="false"), # 使用状況に基づいてスペースを並び替えないようにする
#             dict(domain="com.apple.CrashReporter", key="DialogType", type="string", value="none"), # クラッシュリポーターダイアログを表示しない
#             dict(domain="com.apple.helpviewer", key="DevMode", type="bool", value="true"), # ヘルプを常時前面表示しない
#             dict(domain="com.apple.driver.AppleBluetoothMultitouch.trackpad", key="Clicking", type="int", value="1"), # トラックパッドをタップ = 常時クリック
#             dict(domain="com.apple.LaunchServices", key="LSQuarantine", type="bool", value="false"), # ... 開いてもよろしいですかを表示しない
#         ]
#         for d in darwin_defaults:
#             command("defaults", ["write", d["domain"], d["key"], "-"+d["type"], d["value"]])
# 
# create_directory("$HOME/dev")
# 
# symlink(join_filepath([current_dir, "config", "alacritty"]), "~/.config/alacritty")
# symlink(join_filepath([current_dir, "config", "clojure"]), "~/.config/clojure")
# symlink(join_filepath([current_dir, "config", "vim", "init.vim"]), "~/.vim/vimrc")
# symlink(join_filepath([current_dir, "config", "git"]), "~/.config/git")
# symlink(join_filepath([current_dir, "config", "tig"]), "~/.config/tig")
# 
# command("curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y")
# # git_clone(url="https://github.com/tmux-plugins/tpm", path="~/.local.share.tpm")
# # git_clone(url="https://github.com/rbenv/rbenv", path="~/.rbenv")
# # git_clone(url="https://github.com/rbenv/ruby-build", path="~/.rbenv/plugins/ruby-build")
# # git_clone(url="https://github.com/nodenv/nodenv", path="~/.nodenv")
# # git_clone(url="https://github.com/nodenv/node-build", path="~/.nodenv/plugins/node-build")
# 
# 
# 
# def vscode_setup():
#     if is_darwin():
#         create_directory(path="~/Library/Application\\ Support/Code/User")
#         symlink(join_filepath([current_dir, "config", "Code", "User", "settings.json"]), "~/Library/Application\\ Support/Code/User/settings.json")
#         symlink(join_filepath([current_dir, "config", "Code", "User", "keybindings.json"]), "~/Library/Application\\ Support/Code/User/settings.json")
#         symlink(join_filepath([current_dir, "config", "Code", "User", "locale.json"]), "~/Library/Application\\ Support/Code/User/settings.json")
#         symlink(join_filepath([current_dir, "config", "Code", "User", "projects.json"]), "~/Library/Application\\ Support/Code/User/settings.json")
# 
#     if is_arch_linux():
#         create_directory(path="~/.config/Code/User")
#         symlink(join_filepath([current_dir, "Code", "User", "settings.json"], "~/.config/Code/User/settings.json"))
#         symlink(join_filepath([current_dir, "Code/User/keybindings.json"]), "~/.config/Code/User/keybindings.json")
#         symlink(join_filepath([current_dir, "Code/User/locale.json"]), "~/.config/Code/User/locale.json")
#         symlink(join_filepath([current_dir, "Code/User/projects.json"]), "~/.config/Code/User/projects.json")
# 
# 
# def docker_setup():
#     create_directory("/etc/docker")
#     
#     if is_arch_linux():
#         copy_file(src=join_filepath([current_dir, "files", "daemon.json"]), dest="/etc/docker/daemon.json", permission=0o644)
#         copy_file(src=join_filepath([current_dir, "files", "subuid"]), dest="/etc/subuid", permission=0o644)
#         copy_file(src=join_filepath([current_dir, "files", "subgid"]), dest="/etc/subgid", permission=0o644)
#     
# 
# symlink(join_filepath([current_dir, "config", "hammerspoon"]), "~/.hammerspoon")
# 
# 
# def fish_setup():
#     fish_set_variable("XDG_CONFIG_HOME", "$HOME/.config")
#     fish_set_variable("XDG_CACHE_HOME", "$HOME/.cache")
#     fish_set_variable("XDG_DATA_HOME", "$HOME/.local/share")
#     fish_set_variable("LESSHISTFILE", "$HOME/.cache/less/history")
#     fish_set_variable("MPLAYER_HOME", "$HOME/.config/mplayer")
#     fish_set_variable("INPUTRC", "$HOME/.config/readline/inputrc")
#     fish_set_variable("EDITOR", "vim")
#     fish_set_variable("GOPATH", "$HOME/dev")
#     fish_set_variable("RBENV_ROOT", "$HOME/.rbenv")
#     fish_set_variable("NODENV_ROOT", "$HOME/.nodenv")
# 
#     fish_set_path([
#         "$GOPATH/bin",
#         "$RBENV_ROOT/bin",
#         "$NODENV_ROOT/bin",
#         "$XDG_DATA_HOME/flutter/bin",
#         "$XDG_DATA_HOME/dart-sdk/bin",
#         "$HOME/.local/share/flutter/bin",
#     ])
# 
#     if not executable("fisher"):
#         command("curl https://raw.githubusercontent.com/jorgebucaran/fisher/main/fisher.fish --create-dirs -sLo ~/.config/fish/functions/fisher.fish")
# 
#     symlink(join_filepath([current_dir, "config", "fish", "config.fish"]), "~/.config/fish/config.fish")
#     symlink(join_filepath([current_dir, "config", "fish", "fishfile"]), "~/.config/fish/fishfile")
# 
# def flutter_setup():
#     if is_arch_linux():
#         command("curl -Lo /tmp/flutter.tar.xz https://storage.googleapis.com/flutter_infra/releases/stable/linux/flutter_linux_1.17.1-stable.tar.xz")
#         command("tar xf /tmp/flutter.tar.xz -C /tmp")
#         command("mv /tmp/flutter ~/.local/share/flutter-1.17.1")
#         symlink("~/.local/share/flutter", "~/.local/share/flutter")
# 
# 
# install_packages()
