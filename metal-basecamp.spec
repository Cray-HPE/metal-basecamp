# Copyright 2021 Hewlett Packard Enterprise Development LP
Name: metal-basecamp
License: MIT License
Summary: Datasource for cloud-init metadata
BuildArch: x86_64
Version: %(cat .version)
Release: %(echo ${BUILD_METADATA})
Source: %{name}-%{version}.tar.bz2
Vendor: Cray Inc.
BuildRequires: coreutils
BuildRequires: sed
BuildRequires: skopeo
BuildRequires: pkgconfig(systemd)
Requires: podman
Requires: podman-cni-config

# helps when installing a program whose unit files makes use of a feature only available in a newer systemd version
# If the program is installed on its own, it will have to make do with the available features
# If a newer systemd package is planned to be installed in the same transaction as the program,
# it can be beneficial to have systemd installed first, so that the features have become available by the time program is installed and restarted
%{?systemd_ordering}

%define imagedir %{_sharedstatedir}/cray/container-images/%{name}

%define current_branch %(echo ${GIT_BRANCH} | sed -e 's,/.*$,,')
# Note: Important for basecamp_tag to be the same as used in runPostBuild.sh
%define basecamp_tag   %(echo ${VERSION})

%define bucket csm-docker/unstable
%if "%{current_branch}" == "main"
%undefine bucket
%define bucket csm-docker/unstable
%endif

%if "%{current_branch}" == "release"
%undefine bucket
%define bucket csm-docker/stable
%endif

# This needs to match what is created for the image
%define basecamp_image artifactory.algol60.net/%{bucket}/%{name}:%{basecamp_tag}

%define basecamp_file  cray-metal-basecamp-%{basecamp_tag}.tar

%description
This RPM installs the daemon file for Basecamp, launched through podman.

%prep
env
%setup -q
echo bucket: %{bucket} tag: %{basecamp_tag} current_branch: %{current_branch}
timeout 15m sh -c 'until skopeo inspect docker://%{basecamp_image}; do sleep 10; done'

%build
sed -e 's,@@basecamp-image@@,%{basecamp_image},g' \
    -e 's,@@basecamp-path@@,%{imagedir}/%{basecamp_file},g' \
    -i init/basecamp-init.sh
skopeo copy docker://%{basecamp_image} docker-archive:%{basecamp_file}

%install
install -D -m 0644 -t %{buildroot}%{_unitdir} init/basecamp.service
install -D -m 0755 -t %{buildroot}%{_sbindir} init/basecamp-init.sh
ln -s %{_sbindir}/service %{buildroot}%{_sbindir}/rcbasecamp
install -D -m 0644 -t %{buildroot}%{imagedir} %{basecamp_file}

%clean
rm -f %{basecamp_file}

# These macros will handle sysv initscripts migration transparently (as long as initscripts and systemd services have similar names)
# These also tell systemd about changed unit files--that systemctl daemon-reload should be invoked
%pre
%service_add_pre basecamp.service

%post
%service_add_post basecamp.service

%preun
%service_del_preun basecamp.service

# During package update, %service_del_postun restarts units
%postun
%service_del_postun basecamp.service

%files
%license LICENSE
%doc README.md
%defattr(-,root,root)
%attr(755, root, root) %{_sbindir}/basecamp-init.sh
%attr(644, root, root) %{_unitdir}/basecamp.service
%{_sbindir}/rcbasecamp
%{imagedir}/%{basecamp_file}
%changelog
