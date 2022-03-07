FROM golang:latest as build

WORKDIR /app
COPY ./main.go ./
COPY ./go.mod ./
COPY ./go.sum ./
RUN go get -u -v
RUN CGO_ENABLED=0 go build -o /vmware-powercli -ldflags="-s -w" main.go

FROM ubuntu:21.10

RUN apt-get update && apt-get install -y wget apt-transport-https software-properties-common && \
  wget -q https://packages.microsoft.com/config/ubuntu/20.04/packages-microsoft-prod.deb &&  dpkg -i packages-microsoft-prod.deb && \
  apt-get update && apt-get install -y powershell

RUN pwsh -c "Set-PSRepository -Name PSGallery -InstallationPolicy Trusted" && \
    pwsh -c "Install-Module -Name PSDesiredStateConfiguration" && \
    pwsh -c "Enable-ExperimentalFeature PSDesiredStateConfiguration.InvokeDscResource" && \
    pwsh -c "\$ProgressPreference = \"SilentlyContinue\"; Install-Module VMware.PowerCLI" && \
    pwsh -c "\$ProgressPreference = \"SilentlyContinue\"; Install-Module VMware.vSphereDSC" && \
    pwsh -c "\$ProgressPreference = \"SilentlyContinue\"; Install-Module VMware.PSDesiredStateConfiguration" && \
    pwsh -c "\$ProgressPreference = \"SilentlyContinue\"; Install-Module VMware.vSphere.SsoAdmin" && \
    pwsh -c "\$ProgressPreference = \"SilentlyContinue\"; Install-Module PowerNSX -RequiredVersion 3.0.1174" && \
    pwsh -c "\$ProgressPreference = \"SilentlyContinue\"; Install-Module PowervRA -RequiredVersion 3.6.0" && \
    pwsh -c "\$ProgressPreference = \"SilentlyContinue\"; Set-PowerCLIConfiguration -Scope AllUsers -ParticipateInCEIP \$false -Confirm:\$false" && \
    pwsh -c "\$ProgressPreference = \"SilentlyContinue\"; Set-PowerCLIConfiguration -InvalidCertificateAction Ignore -Scope AllUsers -Confirm:\$false" && \
    echo '$ProgressPreference = "SilentlyContinue"' > /root/.config/powershell/Microsoft.PowerShell_profile.ps1

# Install SRM https://github.com/vmware/PowerCLI-Example-Scripts/tree/master/Modules/SRM
RUN apt-get install ca-certificates zip -y
RUN wget https://github.com/benmeadowcroft/SRM-Cmdlets/archive/master.zip && ls -la && unzip master.zip
RUN pwsh -c "Import-Module -Name .\SRM-Cmdlets-master\Meadowcroft.Srm.psd1"


COPY --from=build /vmware-powercli /
CMD ["/vmware-powercli"]
