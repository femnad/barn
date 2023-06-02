# barn

A small tool for running something which outputs a list, making a selection based on those and saving the selection. Expects the output to be passed to [fzf](https://github.com/junegunn/fzf), uses [bbolt](https://github.com/etcd-io/bbolt) as the database.

## Example

If you have the following configuration:

```yaml
selectors:
  - action: readdir
    id: venv
    settings:
      on_select: source {{ .FullName }}/bin/activate.fish
    target:
      - ~/.local/share/venv
```

The only selector here is one with the id `venv`. When you run `barn` with no arguments, the contents of `~/.local/share/venv` will be printed. When an arguments is provided a `source` command based on the selection is printed. If you put that in a script, you can use `fzf` to activate Python `virtualenv`s which will be sorted according to their usage frequency:

```fish
if test (count $argv) -eq 1
    set vst $argv[1]
else
    barn choose -i venv | fzf +s | read vst
    or return 1
end

eval (barn choose -i venv $vst)
```
