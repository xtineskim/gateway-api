This table is populated from the conformance reports uploaded by project implementations.

| Features                                  | Kong   | cilium   | envoyproxy   | istio   | kumahq   | nginxinc   | projectcontour   | solo.io   |
|:------------------------------------------|:-------|:---------|:-------------|:--------|:---------|:-----------|:-----------------|:----------|
| HTTPRouteBackendRequestHeaderModification | no     | no       | no           | no      | no       | no         | no               | no        |
| HTTPRouteQueryParamMatching               | yes    | yes      | yes          | yes     | yes      | yes        | yes              | yes       |
| HTTPRouteMethodMatching                   | yes    | yes      | yes          | yes     | yes      | yes        | yes              | yes       |
| HTTPRouteResponseHeaderModification       | yes    | yes      | yes          | yes     | yes      | no         | yes              | yes       |
| HTTPRoutePortRedirect                     | no     | yes      | yes          | yes     | yes      | yes        | yes              | yes       |
| HTTPRouteSchemeRedirect                   | no     | yes      | yes          | yes     | yes      | yes        | yes              | yes       |
| HTTPRoutePathRedirect                     | no     | yes      | yes          | yes     | yes      | no         | yes              | yes       |
| HTTPRouteHostRewrite                      | no     | yes      | yes          | yes     | yes      | yes        | yes              | no        |
| HTTPRoutePathRewrite                      | no     | yes      | yes          | yes     | yes      | yes        | yes              | no        |
| HTTPRouteRequestMirror                    | no     | yes      | yes          | yes     | yes      | no         | yes              | no        |
| HTTPRouteRequestMultipleMirrors           | no     | yes      | yes          | yes     | no       | no         | yes              | no        |
| HTTPRouteRequestTimeout                   | no     | yes      | yes          | yes     | no       | no         | yes              | no        |
| HTTPRouteBackendTimeout                   | yes    | yes      | yes          | yes     | no       | no         | yes              | no        |
| HTTPRouteParentRefPort                    | no     | yes      | no           | no      | no       | no         | no               | no        |
