require 'pact/provider/proxy/tasks'

Pact::ProxyVerificationTask.new :javascript do | task |
 task.pact_url ENV["PACT_FILE"], :pact_helper => './tmp/pact_helper'
 task.provider_base_url ENV["PROVIDER_URL"]
end
