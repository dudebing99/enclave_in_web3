## Nitro Enclaves

### Recommended runtime

- hardware: EC2(c5.xlarge)
- os: amazon linux 2023 with `Enclaves enabled`

### Env initialization

- basic

```bash
dnf install golang docker screen htop -y

systemctl start docker
```

- Nitro CLI

```bash
dnf install aws-nitro-enclaves-cli -y
dnf install aws-nitro-enclaves-cli-devel -y

usermod -aG ne root
usermod -aG docker root

nitro-cli --version

systemctl enable --now nitro-enclaves-allocator.service
systemctl enable --now docker
```

### Setting

Enclave configuration file: `/etc/nitro_enclaves/allocator.yaml`

```yaml
# Enclave configuration file.
#
# Location: /etc/nitro_enclaves/allocator.yaml
#
# How much memory to allocate for enclaves (in MiB).
memory_mib: 4000
#
# How many CPUs to reserve for enclaves.
cpu_count: 2
#
# Alternatively, the exact CPUs to be reserved for the enclave can be explicitly
# configured by using `cpu_pool` (like below), instead of `cpu_count`.
# Note: cpu_count and cpu_pool conflict with each other. Only use exactly one of them.
# Example of reserving CPUs 2, 3, and 6 through 9:
# cpu_pool: 2,3,6-9
```

**Take effect after restart**

```bash
systemctl restart nitro-enclaves-allocator.service
```

### Building

```bash
docker build -t hello-enclave:1.0 ./
```

Then, use the above Docker image tag to build the enclave image (`hello.eif`):

```bash
nitro-cli build-enclave --docker-uri hello-enclave:1.0 --output-file /tmp/hello.eif
```

### Running

Now that we have brand-new enclave image, let's use `nitro-cli` to boot it up:

```bash
nitro-cli run-enclave --eif-path /tmp/hello.eif --cpu-count 2 --memory 128 --debug-mode
```

Note: we are running the enclave in debug mode in order to be able to access
      to its console and see our greeting.

We should now be able to connect to the enclave console and see our greeting:

```bash
nitro-cli console --enclave-id $(nitro-cli describe-enclaves | jq -r ".[0].EnclaveID")
```
### Terminate

```bash
nitro-cli terminate-enclave --enclave-id $(nitro-cli describe-enclaves | jq -r ".[0].EnclaveID")
```
