# Copyright 2020 Cray Inc. All Rights Reserved.
Name: basecamp
License: MIT License
Summary: Datasource for cloud-init metadata
BuildArchitectures: noarch
Version: %(cat .version)
Release: %(echo ${BUILD_METADATA})
Source: %{name}-%{version}.tar.bz2
Vendor: Cray Inc.
Requires: podman
%{?systemd_ordering}

%description
This RPM installs the daemon file for Basecamp, launched through podman.

%prep
%setup -q

%build


%install
install -D -m 0644 init/%{name}.service %{buildroot}%{_unitdir}/%{name}.service
mkdir -pv %{buildroot}%{_sbindir}
install -D -m 0755 init/%{name}-init.sh %{buildroot}%{_sbindir}/%{name}-init.sh
ln -s %{_sbindir}/service %{buildroot}%{_sbindir}/rc%{name}

%clean

%pre
%service_add_pre %{name}.service

%post
%service_add_post %{name}.service

%preun
%service_del_preun %{name}.service

%postun
%service_del_postun %{name}.service

%files
%license LICENSE
%doc README.md
%defattr(-,root,root)
%{_unitdir}/%{name}.service
%{_sbindir}/%{name}-init.sh
%{_sbindir}/rc%{name}

%changelog
