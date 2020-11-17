#!/usr/bin/env python3
"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

"""
Pre-run script for services to generate a nghttpx config from a jinja template
and the config/mconfig for the service.
"""

import logging
import os
import socket

from create_oai_certs import generate_mme_certs
from generate_service_config import generate_template_config
from lte.protos.mconfig.mconfigs_pb2 import MME
from magma.common.misc_utils import get_ip_from_if, get_ip_from_if_cidr
from magma.configuration.mconfig_managers import load_service_mconfig
from magma.configuration.service_configs import get_service_config_value

CONFIG_OVERRIDE_DIR = "/var/opt/magma/tmp"
DEFAULT_DNS_IP_PRIMARY_ADDR = "8.8.8.8"
DEFAULT_DNS_IP_SECONDARY_ADDR = "8.8.4.4"
DEFAULT_DNS_IPV6_ADDR = "2001:4860:4860:0:0:0:0:8888"
DEFAULT_P_CSCF_IPV4_ADDR = "172.27.23.150"
DEFAULT_P_CSCF_IPV6_ADDR = "2a12:577:9941:f99c:0002:0001:c731:f114"


def _get_iface_ip(service, iface_config):
    """
    Get the interface IP given its name.
    """
    iface_name = get_service_config_value(service, iface_config, "")
    return get_ip_from_if_cidr(iface_name)


def _get_primary_dns_ip(service_config, iface_config):
    """
    Get dnsd interface IP without netmask.
    If caching is enabled, use the ip of interface that dnsd listens over.
    Otherwise, use dns server from service mconfig.
    """
    if service_config.enable_dns_caching:
        iface_name = get_service_config_value("dnsd", iface_config, "")
        return get_ip_from_if(iface_name)
    else:
        return service_config.dns_primary or DEFAULT_DNS_IP_PRIMARY_ADDR


def _get_secondary_dns_ip(service_config):
    """
    Get the secondary dns ip from the service mconfig.
    """
    return service_config.dns_secondary or DEFAULT_DNS_IP_SECONDARY_ADDR


def _get_ipv4_pcscf_ip(service_config):
    return service_config.ipv4_p_cscf_address or DEFAULT_P_CSCF_IPV4_ADDR


def _get_ipv6_pcscf_ip(service_config):
    return service_config.ipv6_p_cscf_address or DEFAULT_P_CSCF_IPV6_ADDR


def _get_ipv6_dns_ip(service_config):
    """
    Get IPV6 DNS server IP address from service mconfig
    """
    return service_config.ipv6_dns_address or DEFAULT_DNS_IPV6_ADDR


def _get_oai_log_level():
    """
    Convert the logLevel in config into the level which OAI code
    uses. We use OAI's 'TRACE' as the debugging log level and 'CRITICAL'
    as the fatal log level.
    """
    oai_log_level = get_service_config_value("mme", "log_level", "INFO")
    # Translate common log levels to OAI levels
    if oai_log_level == "DEBUG":
        oai_log_level = "TRACE"
    if oai_log_level == "FATAL":
        oai_log_level = "CRITICAL"
    return oai_log_level


def _get_relay_enabled(service_config):
    if service_config.relay_enabled:
        return "yes"
    return "no"


def _get_non_eps_service_control(service_config):
    non_eps_service_control = service_config.non_eps_service_control
    if non_eps_service_control:
        if non_eps_service_control == 0:
            return "OFF"
        elif non_eps_service_control == 1:
            return "CSFB_SMS"
        elif non_eps_service_control == 2:
            return "SMS"
        elif non_eps_service_control == 3:
            return "SMS_ORC8R"
    return "OFF"


def _get_lac(service_config):
    lac = service_config.lac
    if lac:
        return lac
    return 0


def _get_csfb_mcc(service_config):
    csfb_mcc = service_config.csfb_mcc
    if csfb_mcc:
        return csfb_mcc
    return ""


def _get_csfb_mnc(service_config):
    csfb_mnc = service_config.csfb_mnc
    if csfb_mnc:
        return csfb_mnc
    return ""


def _get_identity():
    realm = get_service_config_value("mme", "realm", "")
    return "{}.{}".format(socket.gethostname(), realm)


def _get_enable_nat(service_config):
    nat_enabled = get_service_config_value("mme", "enable_nat", None)

    if nat_enabled is None:
        nat_enabled = service_config.nat_enabled

    if nat_enabled is not None:
        return nat_enabled

    return True


def _get_attached_enodeb_tacs():
    mme_config = load_service_mconfig("mme", MME())
    # attachedEnodebTacs overrides 'tac', which is being deprecated, but for
    # now, both are supported
    tac = mme_config.tac
    attached_enodeb_tacs = mme_config.attached_enodeb_tacs
    if len(attached_enodeb_tacs) == 0:
        return [tac]
    return attached_enodeb_tacs


def _get_apn_correction_map_list(service_config):
    if len(service_config.apn_correction_map_list) == 0:
        return get_service_config_value("mme", "apn_correction_map_list", None)
    return service_config.apn_correction_map_list


def _get_context():
    """
    Create the context which has the interface IP and the OAI log level to use.
    """
    mme_service_config = load_service_mconfig("mme", MME())
    context = {}

    context["mme_s11_ip"] = _get_iface_ip("mme", "s11_iface_name")
    context["sgw_s11_ip"] = _get_iface_ip("spgw", "s11_iface_name")
    context["remote_sgw_ip"] = get_service_config_value("mme", "remote_sgw_ip", "")
    context["s1ap_ip"] = _get_iface_ip("mme", "s1ap_iface_name")
    context["s1u_ip"] = _get_iface_ip("spgw", "s1u_iface_name")
    context["oai_log_level"] = _get_oai_log_level()
    context["ipv4_dns"] = _get_primary_dns_ip(mme_service_config, 'dns_iface_name')
    context["ipv4_sec_dns"] = _get_secondary_dns_ip(mme_service_config)
    context["ipv4_p_cscf_address"] = _get_ipv4_pcscf_ip(mme_service_config)
    context["ipv6_dns"] = _get_ipv6_dns_ip(mme_service_config)
    context["ipv6_p_cscf_address"] = _get_ipv6_pcscf_ip(mme_service_config)
    context["identity"] = _get_identity()
    context["relay_enabled"] = _get_relay_enabled(mme_service_config)
    context["non_eps_service_control"] = _get_non_eps_service_control(mme_service_config)
    context["csfb_mcc"] = _get_csfb_mcc(mme_service_config)
    context["csfb_mnc"] = _get_csfb_mnc(mme_service_config)
    context["lac"] = _get_lac(mme_service_config)
    context["use_stateless"] = get_service_config_value("mme", "use_stateless", "")
    context["attached_enodeb_tacs"] = _get_attached_enodeb_tacs()
    context["enable_nat"] = _get_enable_nat(mme_service_config)
    # set ovs params
    for key in (
        "ovs_bridge_name",
        "ovs_gtp_port_number",
        "ovs_mtr_port_number",
        "ovs_internal_sampling_port_number",
        "ovs_internal_sampling_fwd_tbl",
        "ovs_uplink_port_number",
        "ovs_uplink_mac",
    ):
        context[key] = get_service_config_value("spgw", key, "")
    context["enable_apn_correction"] = get_service_config_value("mme", "enable_apn_correction", "")
    context["apn_correction_map_list"] = _get_apn_correction_map_list(mme_service_config)
    return context


def main():
    logging.basicConfig(
        level=logging.INFO, format="[%(asctime)s %(levelname)s %(name)s] %(message)s"
    )
    context = _get_context()
    generate_template_config("spgw", "spgw", CONFIG_OVERRIDE_DIR, context.copy())
    generate_template_config("mme", "mme", CONFIG_OVERRIDE_DIR, context.copy())
    generate_template_config("mme", "mme_fd", CONFIG_OVERRIDE_DIR, context.copy())
    cert_dir = get_service_config_value("mme", "cert_dir", "")
    generate_mme_certs(os.path.join(cert_dir, "freeDiameter"))


if __name__ == "__main__":
    main()
