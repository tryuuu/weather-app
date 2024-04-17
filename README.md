## 実行
`chmod +x exec.sh`

`./exec.sh`

トラフィックの確認

`istioctl dashboard kiali`

`for i in {1..20}; do curl http://127.0.0.1:56357; done`

デバッグ

`istioctl analyze -n development`

`istioctl analyze -n istio-system`

`kubectl get virtualservice -n development -o yaml`

## ログイン
user: tryu

pass: k8s
