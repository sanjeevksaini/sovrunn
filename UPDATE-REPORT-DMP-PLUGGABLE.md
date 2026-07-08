# Update Report: DMP as Pluggable Management Plane

Applied repository-wide architecture consistency update.

## Updated files

- docs/architecture/control-plane/control-plane.md
- docs/architecture/control-plane/control-plane-map.md
- docs/architecture/control-plane/management-plane.md
- docs/architecture/control-plane/datastore-management-plane/datastore-management-plane.md
- docs/architecture/control-plane/management-plane-framework/management-plane-framework.md
- docs/architecture/control-plane/management-plane-framework/management-plane-registry.md
- docs/architecture/control-plane/management-plane-framework/management-plane-manifest.md
- docs/architecture/control-plane/management-plane-framework/management-plane-controller-runtime.md
- docs/architecture/control-plane/management-plane-framework/management-plane-admission.md
- docs/architecture/control-plane/management-plane-framework/management-plane-conformance.md
- docs/foundation/glossary.md
- mkdocs.yml

## Summary

- DMP is now consistently modeled as a pluggable management plane inside SDE Control Plane.
- Added Management Plane Framework architecture documents.
- Clarified `sde-dmp-controller` as DMP Controller Runtime, not the whole DMP.
- Updated glossary and MkDocs navigation.
