################################################################################

%global crc_check pushd ../SOURCES ; sha512sum -c %{SOURCE100} ; popd

################################################################################

%define debug_package  %{nil}

################################################################################

Summary:        Tool for encrypting/decrypting arbitrary data streams
Name:           siocrypt
Version:        0.1.0
Release:        0%{?dist}
Group:          Applications/System
License:        Apache License, Version 2.0
URL:            https://kaos.sh/siocrypt

Source0:        https://source.kaos.st/%{name}/%{name}-%{version}.tar.bz2

Source100:      checksum.sha512

BuildRoot:      %{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)

BuildRequires:  golang >= 1.23

Provides:       %{name} = %{version}-%{release}

################################################################################

%description
siocrypt is a tool for encrypting/decrypting arbitrary data streams using Data
At Rest Encryption (DARE).

################################################################################

%prep
%{crc_check}

%setup -q
if [[ ! -d "%{name}/vendor" ]] ; then
  echo -e "----\nThis package requires vendored dependencies\n----"
  exit 1
elif [[ -f "%{name}/%{name}" ]] ; then
  echo -e "----\nSources must not contain precompiled binaries\n----"
  exit 1
fi

%build
pushd %{name}
  %{__make} %{?_smp_mflags} all
  cp LICENSE ..
popd

%install
rm -rf %{buildroot}

install -dDm 755 %{buildroot}%{_bindir}

install -pm 755 %{name}/%{name} \
                %{buildroot}%{_bindir}/

install -dDm 755 %{buildroot}%{_mandir}/man1

./%{name}/%{name} --generate-man > %{buildroot}%{_mandir}/man1/%{name}.1

install -dDm 755 %{buildroot}%{_sysconfdir}/bash_completion.d
install -dDm 755 %{buildroot}%{_datadir}/zsh/site-functions
install -dDm 755 %{buildroot}%{_datarootdir}/fish/vendor_completions.d

./%{name}/%{name} --completion=bash 1> %{buildroot}%{_sysconfdir}/bash_completion.d/%{name}
./%{name}/%{name} --completion=zsh 1> %{buildroot}%{_datadir}/zsh/site-functions/_%{name}
./%{name}/%{name} --completion=fish 1> %{buildroot}%{_datarootdir}/fish/vendor_completions.d/%{name}.fish

%clean
rm -rf %{buildroot}

################################################################################

%files
%defattr(-,root,root,-)
%doc LICENSE
%{_bindir}/%{name}
%{_mandir}/man1/%{name}.1.*
%{_sysconfdir}/bash_completion.d/%{name}
%{_datadir}/zsh/site-functions/_%{name}
%{_datarootdir}/fish/vendor_completions.d/%{name}.fish

################################################################################

%changelog
* Fri Jun 27 2025 Anton Novojilov <andy@essentialkaos.com> - 0.0.3-0
- Initial build for kaos-repo
