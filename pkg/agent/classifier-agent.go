// Generated by *go generate* - DO NOT EDIT
/*
Copyright 2022. projectsveltos.io. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package agent

var classifierAgentYAML = []byte(`apiVersion: v1
kind: Namespace
metadata:
  labels:
    control-plane: classifier-agent
  name: projectsveltos
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: classifier-agent-manager
  namespace: projectsveltos
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: classifier-agent-leader-election-role
  namespace: projectsveltos
rules:
- apiGroups:
  - ""
  resources:
  - configmaps
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
- apiGroups:
  - coordination.k8s.io
  resources:
  - leases
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
- apiGroups:
  - ""
  resources:
  - events
  verbs:
  - create
  - patch
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  creationTimestamp: null
  name: classifier-agent-manager-role
rules:
- apiGroups:
  - '*'
  resources:
  - '*'
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - lib.projectsveltos.io
  resources:
  - classifierreports
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
- apiGroups:
  - lib.projectsveltos.io
  resources:
  - classifierreports/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - lib.projectsveltos.io
  resources:
  - classifiers
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - lib.projectsveltos.io
  resources:
  - classifiers/finalizers
  verbs:
  - update
- apiGroups:
  - lib.projectsveltos.io
  resources:
  - classifiers/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - projectsveltos.io
  resources:
  - nodes
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - projectsveltos.io
  resources:
  - nodes/finalizers
  verbs:
  - update
- apiGroups:
  - projectsveltos.io
  resources:
  - nodes/status
  verbs:
  - get
  - patch
  - update
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: classifier-agent-metrics-reader
rules:
- nonResourceURLs:
  - /metrics
  verbs:
  - get
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: classifier-agent-proxy-role
rules:
- apiGroups:
  - authentication.k8s.io
  resources:
  - tokenreviews
  verbs:
  - create
- apiGroups:
  - authorization.k8s.io
  resources:
  - subjectaccessreviews
  verbs:
  - create
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: classifier-agent-leader-election-rolebinding
  namespace: projectsveltos
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: classifier-agent-leader-election-role
subjects:
- kind: ServiceAccount
  name: classifier-agent-manager
  namespace: projectsveltos
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: classifier-agent-manager-rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: classifier-agent-manager-role
subjects:
- kind: ServiceAccount
  name: classifier-agent-manager
  namespace: projectsveltos
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: classifier-agent-proxy-rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: classifier-agent-proxy-role
subjects:
- kind: ServiceAccount
  name: classifier-agent-manager
  namespace: projectsveltos
---
apiVersion: v1
data:
  controller_manager_config.yaml: |
    apiVersion: controller-runtime.sigs.k8s.io/v1alpha1
    kind: ControllerManagerConfig
    health:
      healthProbeBindAddress: :8081
    metrics:
      bindAddress: 127.0.0.1:8080
    webhook:
      port: 9443
    leaderElection:
      leaderElect: true
      resourceName: e8afc439.projectsveltos.io
    # leaderElectionReleaseOnCancel defines if the leader should step down volume
    # when the Manager ends. This requires the binary to immediately end when the
    # Manager is stopped, otherwise, this setting is unsafe. Setting this significantly
    # speeds up voluntary leader transitions as the new leader don't have to wait
    # LeaseDuration time first.
    # In the default scaffold provided, the program ends immediately after
    # the manager stops, so would be fine to enable this option. However,
    # if you are doing or is intended to do any operation such as perform cleanups
    # after the manager stops then its usage might be unsafe.
    # leaderElectionReleaseOnCancel: true
kind: ConfigMap
metadata:
  name: classifier-agent-manager-config
  namespace: projectsveltos
---
apiVersion: v1
kind: Service
metadata:
  labels:
    control-plane: classifier-agent
  name: classifier-agent-manager-metrics-service
  namespace: projectsveltos
spec:
  ports:
  - name: https
    port: 8443
    protocol: TCP
    targetPort: https
  selector:
    control-plane: classifier-agent
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    control-plane: classifier-agent
  name: classifier-agent-manager
  namespace: projectsveltos
spec:
  replicas: 1
  selector:
    matchLabels:
      control-plane: classifier-agent
  template:
    metadata:
      annotations:
        kubectl.kubernetes.io/default-container: manager
      labels:
        control-plane: classifier-agent
    spec:
      containers:
      - args:
        - --health-probe-bind-address=:8081
        - --metrics-bind-address=127.0.0.1:8080
        - --leader-elect
        - --v=5
        command:
        - /manager
        image: gianlucam76/classifier-agent-manager-amd64:v0.2.0
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8081
          initialDelaySeconds: 15
          periodSeconds: 20
        name: manager
        readinessProbe:
          httpGet:
            path: /readyz
            port: 8081
          initialDelaySeconds: 5
          periodSeconds: 10
        resources:
          limits:
            cpu: 500m
            memory: 128Mi
          requests:
            cpu: 10m
            memory: 64Mi
        securityContext:
          allowPrivilegeEscalation: false
          capabilities:
            drop:
            - ALL
      - args:
        - --secure-listen-address=0.0.0.0:8443
        - --upstream=http://127.0.0.1:8080/
        - --logtostderr=true
        - --v=0
        image: gcr.io/kubebuilder/kube-rbac-proxy:v0.12.0
        name: kube-rbac-proxy
        ports:
        - containerPort: 8443
          name: https
          protocol: TCP
        resources:
          limits:
            cpu: 500m
            memory: 128Mi
          requests:
            cpu: 5m
            memory: 64Mi
        securityContext:
          allowPrivilegeEscalation: false
          capabilities:
            drop:
            - ALL
      securityContext:
        runAsNonRoot: true
      serviceAccountName: classifier-agent-manager
      terminationGracePeriodSeconds: 10
`)
