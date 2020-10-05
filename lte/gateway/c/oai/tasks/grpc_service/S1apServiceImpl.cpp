/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */
#include "S1apServiceImpl.h"

#include <string>

extern "C" {
#include "s1ap_state.h"
#include "log.h"
#include "hashtable.h"
}

namespace grpc {
class ServerContext;
}  // namespace grpc

namespace magma {
namespace orc8r {
class Void;
}  // namespace orc8r
} // namespace magma

using grpc::ServerContext;
using grpc::Status;
using magma::EnbConnectedResult;
using magma::S1apService;

namespace magma {
using namespace lte;
using namespace orc8r;

S1apServiceImpl::S1apServiceImpl() {}

Status S1apServiceImpl::GetEnbConnected(
    ServerContext* context, const Void* request,
    EnbConnectedResult* response) {
  OAILOG_DEBUG(LOG_UTIL, "Received EnbConnected GRPC request\n");

  s1ap_state_t* s1ap_state = get_s1ap_state(false);
  if (s1ap_state != nullptr) {
    hashtable_rc_t ht_rc;
    hashtable_key_array_t* ht_keys = hashtable_ts_get_keys(&s1ap_state->enbs);
    if (ht_keys == nullptr) {
      return Status::OK;
    }

    for (uint32_t i = 0; i < ht_keys->num_keys; i++) {
      enb_description_t* enb_ref;
      ht_rc = hashtable_ts_get(
          &s1ap_state->enbs, (hash_key_t) ht_keys->keys[i], (void**) &enb_ref);
      if (ht_rc == HASH_TABLE_OK) {
        response->add_enb_ids(enb_ref->enb_id);
      }
    }
  }

  return Status::OK;
}

}  // namespace magma

