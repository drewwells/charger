APPDIR=.
APPYAML=$(APPDIR)/app.yaml

YAWN:=$(shell which technicolor-yawn || echo cat)

serve:
	dev_appserver.py $(APPYAML) 2>&1 | $(YAWN)

deploy:
	appcfg.py update_indexes $(APPDIR) --noauth_local_webserver
	gcloud app deploy
