FROM golang:latest as build

WORKDIR /app
COPY ./main.go ./
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
RUN CGO_ENABLED=0 go build -o /vmware-powercli -ldflags="-s -w" main.go

FROM mcr.microsoft.com/powershell:ubuntu-20.04

ENV TERM linux

WORKDIR /root

RUN apt-get update && apt-get install curl wget git zip tar ca-certificates unzip -y

RUN echo "/usr/bin/pwsh" >> /etc/shells && \
    echo "/bin/pwsh" >> /etc/shells && \
    pwsh -c "Set-PSRepository -Name PSGallery -InstallationPolicy Trusted" && \
    pwsh -c "Install-Module -Name PSDesiredStateConfiguration" && \
    pwsh -c "Enable-ExperimentalFeature PSDesiredStateConfiguration.InvokeDscResource" && \
    pwsh -c "\$ProgressPreference = \"SilentlyContinue\"; Install-Module VMware.PowerCLI" && \
    pwsh -c "\$ProgressPreference = \"SilentlyContinue\"; Install-Module VMware.vSphereDSC" && \
    pwsh -c "\$ProgressPreference = \"SilentlyContinue\"; Install-Module VMware.PSDesiredStateConfiguration" && \
    pwsh -c "\$ProgressPreference = \"SilentlyContinue\"; Install-Module VMware.vSphere.SsoAdmin" && \
    pwsh -c "\$ProgressPreference = \"SilentlyContinue\"; Install-Module PowerNSX -RequiredVersion 3.0.1174" && \
    pwsh -c "\$ProgressPreference = \"SilentlyContinue\"; Install-Module PowervRA -RequiredVersion 3.6.0" && \
    pwsh -c "\$ProgressPreference = \"SilentlyContinue\"; Set-PowerCLIConfiguration -Scope AllUsers -ParticipateInCEIP \$false -Confirm:\$false" && \
    pwsh -c "\$ProgressPreference = \"SilentlyContinue\"; Set-PowerCLIConfiguration -InvalidCertificateAction Ignore -Scope AllUsers -Confirm:\$false" 

RUN curl -o ./PowerCLI-Example-Scripts.zip -J -L https://github.com/vmware/PowerCLI-Example-Scripts/archive/03272c1d2db26a525b31c930e3bf3d20d34468e0.zip && \
    unzip PowerCLI-Example-Scripts.zip && \
    rm -f PowerCLI-Example-Scripts.zip && \
    mv ./PowerCLI-Example-Scripts-* ./PowerCLI-Example-Scripts && \
    mv ./PowerCLI-Example-Scripts/Modules/* /opt/microsoft/powershell/7/Modules

RUN wget https://github.com/benmeadowcroft/SRM-Cmdlets/archive/master.zip && unzip master.zip && rm master.zip
RUN pwsh -c "Import-Module -Name .\SRM-Cmdlets-master\Meadowcroft.Srm.psd1"
RUN echo '$ProgressPreference = "SilentlyContinue"' > /root/.config/powershell/Microsoft.PowerShell_profile.ps1  
 
COPY --from=build /vmware-powercli /

CMD ["/vmware-powercli"]
