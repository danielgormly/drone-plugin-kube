package main

import (
	"log"
)

type (
	Repo struct {
		Owner string
		Name  string
	}

	Build struct {
		Tag     string
		Event   string
		Number  int
		Commit  string
		Ref     string
		Branch  string
		Author  string
		Status  string
		Link    string
		Started int64
		Created int64
	}

	Job struct {
		Started int64
	}

	Config struct {
		Ca        string
		Server    string
		Token     string
		Namespace string
		Template  string
	}

	Plugin struct {
		Repo   Repo
		Build  Build
		Config Config
		Job    Job
	}
)

func (p Plugin) Exec() error {

	if p.Config.Server == "" {
		log.Fatal("KUBE_SERVER is not defined")
	}
	if p.Config.Token == "" {
		log.Fatal("KUBE_TOKEN is not defined")
	}
	if p.Config.Ca == "" {
		log.Fatal("KUBE_CA is not defined")
	}
	if p.Config.Namespace == "" {
		p.Config.Namespace = "default"
	}
	if p.Config.Template == "" {
		log.Fatal("KUBE_TEMPLATE, or template must be defined")
	}

	// // connect to Kubernetes
	// clientset, err := p.createKubeClient()
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// parse the template file and do substitutions
	txt, err := openAndSub(p.Config.Template, p)
	if err != nil {
		return err
	}
	// // convert txt back to []byte and convert to json
	// json, err := utilyaml.ToJSON([]byte(txt))
	// if err != nil {
	// 	return err
	// }

	// var dep v1beta1.Deployment

	// e := runtime.DecodeInto(api.Codecs.UniversalDecoder(), json, &dep)
	// if e != nil {
	// 	log.Fatal("Error decoding yaml file to json", e)
	// }
	// // check and see if there is a deployment already.  If there is, update it.
	// oldDep, err := findDeployment(dep.ObjectMeta.Name, dep.ObjectMeta.Namespace, clientset)
	// if err != nil {
	// 	return err
	// }
	// if oldDep.ObjectMeta.Name == dep.ObjectMeta.Name {
	// 	// update the existing deployment, ignore the deployment that it comes back with
	// 	_, err = clientset.ExtensionsV1beta1().Deployments(p.Config.Namespace).Update(&dep)
	// 	return err
	// }
	// // create the new deployment since this never existed.
	// _, err = clientset.ExtensionsV1beta1().Deployments(p.Config.Namespace).Create(&dep)

	// return err
}
