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
Requires: podman
Requires: podman-cni-config
Provides: basecamp
%{?systemd_ordering}

%define imagedir %{_sharedstatedir}/cray/container-images/%{name}

%define current_branch %(echo ${GIT_BRANCH} | sed -e 's,/.*$,,')
# Note: Important for basecamp_tag to be the same as used in runPostBuild.sh
%define basecamp_tag   %{version}-%(git rev-parse --short HEAD)

%define bucket csm-docker-unstable-local
%if "%{current_branch}" == "main"
%undefine bucket
%define bucket csm-docker-master-local
%endif

%if "%{current_branch}" == "release"
%undefine bucket
%define bucket csm-docker-stable-local
%endif

%define basecamp_image arti.dev.cray.com/%{bucket}/metal-basecamp:%{basecamp_tag}
%define basecamp_file  metal-basecamp-%{basecamp_tag}.tar

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

%pre
%service_add_pre basecamp.service

%post
%service_add_post basecamp.service

%preun
%service_del_preun basecamp.service

%postun
%service_del_postun basecamp.service

%files
%license LICENSE
%doc README.md
%defattr(-,root,root)
%{_unitdir}/basecamp.service
%{_sbindir}/basecamp-init.sh
%{_sbindir}/rcbasecamp
%{imagedir}/%{basecamp_file}

%changelog
