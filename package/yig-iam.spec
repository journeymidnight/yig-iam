%global debug_package %{nil}
%global __strip /bin/true

Name:           yig-iam
Version:        %{ver}
Release:        %{rel}%{?dist}

Summary:	A simple IAM service for yig

Group:		SDS
License:	GPL
URL:		http://github.com/journeymidnight
Source0:	%{name}-%{version}-%{rel}.tar.gz
BuildRoot:	%(mktemp -ud %{_tmppath}/%{name}-%{version}-%{release}-XXXXXX)
#BuildRequires:  


%description


%prep
%setup -q -n %{name}-%{version}-%{rel}


%build
make

%install
rm -rf %{buildroot}
install -D -m 755 yig-iam %{buildroot}%{_bindir}/yig-iam
install -D -m 755 yig-iam-tools %{buildroot}%{_bindir}/yig-iam-tools
install -D -m 644 package/yig-iam.logrotate %{buildroot}/etc/logrotate.d/yig-iam.logrotate
install -D -m 644 package/yig-iam.service   %{buildroot}/usr/lib/systemd/system/yig-iam.service
install -D -m 644 conf/conf.toml %{buildroot}%{_sysconfdir}/yig-iam/conf.toml
install -D -m 644 conf/basic_model.conf %{buildroot}%{_sysconfdir}/yig-iam/basic_model.conf
install -D -m 644 conf/casbin.conf %{buildroot}%{_sysconfdir}/yig-iam/casbin.conf
install -d %{buildroot}/var/log/yig-iam/


%post
systemctl enable yig-iam


%preun

%clean
rm -rf %{buildroot}

%files
%defattr(-,root,root,-)
%dir /var/log/yig-iam/
%dir /etc/yig-iam/
%config(noreplace) /etc/yig-iam/conf.toml
%config(noreplace) /etc/yig-iam/basic_model.conf
%config(noreplace) /etc/yig-iam/casbin.conf
/usr/bin/yig-iam
/usr/bin/yig-iam-tools
/etc/logrotate.d/yig-iam.logrotate
/usr/lib/systemd/system/yig-iam.service


%changelog
